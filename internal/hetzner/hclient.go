package hetzner

import (
	"bytes"
	"context"
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
		var token []byte

		if config.HetznerTokenSecret.Name != "" {
			secret, err := kubeClient.CoreV1().Secrets(namespace).Get(
				ctx,
				config.HetznerTokenSecret.Name,
				v1.GetOptions{},
			)
			if err != nil {
				return nil, err
			}

			var ok bool
			token, ok = secret.Data[config.HetznerTokenSecret.Key]
			if !ok {
				return nil, fmt.Errorf(
					"secret %s in namespace %s does not contain key %s",
					config.HetznerTokenSecret.Name,
					namespace,
					config.HetznerTokenSecret.Key,
				)
			}
		} else if config.HetznerTokenFilePath != "" {
			data, err := os.ReadFile(config.HetznerTokenFilePath)
			if err != nil {
				return nil, fmt.Errorf("error reading hetzner token file: %w", err)
			}
			token = data
		}

		token = bytes.TrimSpace(token)
		if len(token) == 0 {
			return nil, fmt.Errorf("hetzner token not provided (set tokenSecretKeyRef or tokenFilePath)")
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

func MockHClientBuilder(client *hcloud.Client) HClientBuilderFunc {
	return func(
		_ context.Context,
		_ string,
		_ Config,
	) (*hcloud.Client, error) {
		return client, nil
	}
}
