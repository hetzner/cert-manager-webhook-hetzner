package hetzner

import (
	"context"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestBuildClient(t *testing.T) {
	t.Run("Registry", func(t *testing.T) {
		registry := prometheus.NewRegistry()

		namespace := "default"

		secret := &v1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "hcloud",
				Namespace: namespace,
			},
			Data: map[string][]byte{
				"token": []byte("test-token"),
			},
		}

		kubeClientSet := fake.NewClientset(secret)

		config := Config{
			HetznerTokenSecret: SecretKeyRef{
				Name: "hcloud",
				Key:  "token",
			},
		}

		builder := NewHClientBuilder(kubeClientSet, registry)
		client, err := builder(t.Context(), namespace, config)
		require.NoError(t, err)

		_, _ = client.Location.All(context.Background())

		require.NoError(t, testutil.GatherAndCompare(registry,
			strings.NewReader(`
# HELP hcloud_api_requests_total A counter for requests to the hcloud api per endpoint.
# TYPE hcloud_api_requests_total counter
hcloud_api_requests_total{api_endpoint="/locations",code="401",method="get"} 1
`),
			"hcloud_api_requests_total",
		))
	})

	t.Run("SecretMissingKey", func(t *testing.T) {
		namespace := "default"

		secret := &v1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "hcloud",
				Namespace: namespace,
			},
			Data: map[string][]byte{
				"other-key": []byte("test-token"),
			},
		}

		kubeClientSet := fake.NewClientset(secret)

		config := Config{
			HetznerTokenSecret: SecretKeyRef{
				Name: "hcloud",
				Key:  "token",
			},
		}

		builder := NewHClientBuilder(kubeClientSet, prometheus.NewRegistry())
		_, err := builder(t.Context(), namespace, config)
		require.EqualError(t, err, "secret hcloud in namespace default does not contain key token")
	})

	t.Run("TokenFromFile", func(t *testing.T) {
		tokenPath := filepath.Join(t.TempDir(), "token")
		require.NoError(t, os.WriteFile(tokenPath, []byte("file-token\n"), 0o600))

		config := Config{
			HetznerTokenFilePath: tokenPath,
		}

		builder := NewHClientBuilder(fake.NewClientset(), prometheus.NewRegistry())
		client, err := builder(t.Context(), "default", config)
		require.NoError(t, err)
		require.NotNil(t, client)
	})

	t.Run("TokenFileMissing", func(t *testing.T) {
		tokenPath := filepath.Join(t.TempDir(), "does-not-exist")
		config := Config{
			HetznerTokenFilePath: tokenPath,
		}

		builder := NewHClientBuilder(fake.NewClientset(), prometheus.NewRegistry())
		_, err := builder(t.Context(), "default", config)
		require.ErrorContains(t, err, "error reading hetzner token file")
		require.ErrorIs(t, err, fs.ErrNotExist)
	})

	t.Run("TokenFileEmpty", func(t *testing.T) {
		tokenPath := filepath.Join(t.TempDir(), "token")
		require.NoError(t, os.WriteFile(tokenPath, []byte("   \n"), 0o600))

		config := Config{
			HetznerTokenFilePath: tokenPath,
		}

		builder := NewHClientBuilder(fake.NewClientset(), prometheus.NewRegistry())
		_, err := builder(t.Context(), "default", config)
		require.EqualError(t, err, "hetzner token is empty")
	})

	t.Run("NoTokenConfigured", func(t *testing.T) {
		builder := NewHClientBuilder(fake.NewClientset(), prometheus.NewRegistry())
		_, err := builder(t.Context(), "default", Config{})
		require.EqualError(t, err, "hetzner token not provided (set tokenSecretKeyRef or tokenFilePath)")
	})

	t.Run("BothTokenSourcesConfigured", func(t *testing.T) {
		config := Config{
			HetznerTokenSecret:   SecretKeyRef{Name: "hcloud", Key: "token"},
			HetznerTokenFilePath: "/tmp/hetzner",
		}

		builder := NewHClientBuilder(fake.NewClientset(), prometheus.NewRegistry())
		_, err := builder(t.Context(), "default", config)
		require.EqualError(t, err, "hetzner token provided in both tokenSecretKeyRef and tokenFilePath")
	})

	t.Run("InvalidHTTPTimeout", func(t *testing.T) {
		config := Config{
			HetznerTokenFilePath: writeTestToken(t),
			HTTPTimeout:          "not-a-duration",
		}

		builder := NewHClientBuilder(fake.NewClientset(), prometheus.NewRegistry())
		_, err := builder(t.Context(), "default", config)
		require.ErrorContains(t, err, "error parsing httpTimeout")
	})

	t.Run("CustomHTTPTimeoutIsAppliedToClient", func(t *testing.T) {
		// A server that responds slower than the configured timeout allows. The call
		// context is bounded well below the SDK's first retry backoff (1s), so a
		// correctly-applied short timeout fails fast within that bound; if the
		// configured value were ignored in favor of the (much larger) default, the
		// request would still be in flight when the context deadline cuts it off with
		// a different error instead.
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			time.Sleep(50 * time.Millisecond)
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{}`))
		}))
		t.Cleanup(server.Close)

		config := Config{
			HetznerTokenFilePath: writeTestToken(t),
			HCloudEndpoint:       server.URL,
			HTTPTimeout:          "10ms",
		}

		builder := NewHClientBuilder(fake.NewClientset(), prometheus.NewRegistry())
		client, err := builder(t.Context(), "default", config)
		require.NoError(t, err)

		ctx, cancel := context.WithTimeout(t.Context(), 200*time.Millisecond)
		defer cancel()

		_, err = client.Location.All(ctx)
		require.ErrorContains(t, err, "Client.Timeout")
	})
}

func writeTestToken(t *testing.T) string {
	t.Helper()

	tokenPath := filepath.Join(t.TempDir(), "token")
	require.NoError(t, os.WriteFile(tokenPath, []byte("test-token"), 0o600))
	return tokenPath
}

func TestResolveHTTPTimeout(t *testing.T) {
	t.Run("Empty falls back to default", func(t *testing.T) {
		d, err := resolveHTTPTimeout("")
		require.NoError(t, err)
		assert.Equal(t, DefaultHTTPTimeout, d)
	})

	t.Run("Valid duration is parsed", func(t *testing.T) {
		d, err := resolveHTTPTimeout("45s")
		require.NoError(t, err)
		assert.Equal(t, 45*time.Second, d)
	})

	t.Run("Invalid duration returns an error", func(t *testing.T) {
		_, err := resolveHTTPTimeout("not-a-duration")
		require.ErrorContains(t, err, "error parsing httpTimeout")
	})

	t.Run("Zero duration is rejected", func(t *testing.T) {
		_, err := resolveHTTPTimeout("0s")
		require.ErrorContains(t, err, "httpTimeout must be positive")
	})

	t.Run("Negative duration is rejected", func(t *testing.T) {
		_, err := resolveHTTPTimeout("-5s")
		require.ErrorContains(t, err, "httpTimeout must be positive")
	})
}
