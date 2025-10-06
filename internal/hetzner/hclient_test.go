package hetzner

import (
	"context"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
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
				"token": {},
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
}
