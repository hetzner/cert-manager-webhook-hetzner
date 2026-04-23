package hetzner

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/hetzner/cert-manager-webhook-hetzner/internal/version"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

// HClientBuilderFunc is a function that builds a [hcloud.Client] from a secret stored in kubernetes.
type HClientBuilderFunc func(
	ctx context.Context,
	namespace string,
	config Config,
) (*hcloud.Client, error)

func NewHClientBuilder(kubeClient kubernetes.Interface, registry prometheus.Registerer) HClientBuilderFunc {
	return func(
		ctx context.Context,
		namespace string,
		config Config,
	) (*hcloud.Client, error) {
		var (
			token []byte
			err   error
		)
		secretConfigured := config.HetznerTokenSecret.Name != ""
		fileConfigured := config.HetznerTokenFilePath != ""

		switch {
		case secretConfigured && fileConfigured:
			return nil, errors.New("hetzner token provided in both tokenSecretKeyRef and tokenFilePath")
		case secretConfigured:
			token, err = loadTokenFromSecret(ctx, kubeClient, namespace, config.HetznerTokenSecret)
		case fileConfigured:
			token, err = loadTokenFromFile(config.HetznerTokenFilePath)
		default:
			return nil, errors.New("hetzner token not provided (set tokenSecretKeyRef or tokenFilePath)")
		}
		if err != nil {
			return nil, err
		}

		token = bytes.TrimSpace(token)
		if len(token) == 0 {
			return nil, errors.New("hetzner token is empty")
		}

		clientOpts := []hcloud.ClientOption{
			hcloud.WithToken(string(token)),
			hcloud.WithInstrumentation(registry),
			hcloud.WithApplication("cert-manager-webhook-hetzner", version.Version),
			hcloud.WithHTTPClient(&http.Client{Timeout: 15 * time.Second}),
		}
		if config.HCloudEndpoint != "" {
			clientOpts = append(clientOpts, hcloud.WithEndpoint(config.HCloudEndpoint))
		}

		return hcloud.NewClient(clientOpts...), nil
	}
}

func loadTokenFromSecret(
	ctx context.Context,
	kubeClient kubernetes.Interface,
	namespace string,
	ref SecretKeyRef,
) ([]byte, error) {
	secret, err := kubeClient.CoreV1().Secrets(namespace).Get(ctx, ref.Name, v1.GetOptions{})
	if err != nil {
		return nil, err
	}

	token, ok := secret.Data[ref.Key]
	if !ok {
		return nil, fmt.Errorf(
			"secret %s in namespace %s does not contain key %s",
			ref.Name, namespace, ref.Key,
		)
	}
	return token, nil
}

func loadTokenFromFile(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading hetzner token file: %w", err)
	}
	return data, nil
}

func MockHClientBuilder(client *hcloud.Client) HClientBuilderFunc {
	return func(
		_ context.Context,
		_ string,
		_ Config,
	) (*hcloud.Client, error) {
		return client, nil
	}
}
