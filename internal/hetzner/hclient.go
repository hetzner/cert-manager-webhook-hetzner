package hetzner

import (
	"context"
	"fmt"

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
		secret, err := kubeClient.CoreV1().Secrets(namespace).Get(
			ctx,
			config.HetznerTokenSecret.Name,
			v1.GetOptions{},
		)
		if err != nil {
			return nil, err
		}

		token, ok := secret.Data[config.HetznerTokenSecret.Key]
		if !ok {
			return nil, fmt.Errorf(
				"secret %s in namespace %s does not contain key %s",
				config.HetznerTokenSecret.Name,
				namespace,
				config.HetznerTokenSecret.Key,
			)
		}

		clientOpts := []hcloud.ClientOption{
			hcloud.WithToken(string(token)),
			hcloud.WithInstrumentation(registry),
			hcloud.WithApplication("cert-manager-webhook-hetzner", version.Version),
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
