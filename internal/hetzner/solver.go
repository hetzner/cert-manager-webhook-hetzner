package hetzner

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/cert-manager/cert-manager/pkg/acme/webhook"
	"github.com/cert-manager/cert-manager/pkg/acme/webhook/apis/acme/v1alpha1"
	"github.com/prometheus/client_golang/prometheus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/hetznercloud/hcloud-go/v2/hcloud/exp/zoneutil"
)

// TTL defines the TTL of the records we create. The TTL is small, as these
// records are only needed for the ACME process.
const TTL int = 300

type Solver struct {
	logger   *slog.Logger
	registry prometheus.Registerer

	hClientBuilder HClientBuilderFunc
}

var _ webhook.Solver = (*Solver)(nil)

func New(logger *slog.Logger, registry prometheus.Registerer) *Solver {
	return &Solver{
		logger:   logger,
		registry: registry,
	}
}

// Name is used as the name for this DNS solver when referencing it on the ACME Issuer
// resource.
// This should be unique **within the group name**, i.e. you can have two solvers
// configured with the same Name() **so long as they do not co-exist within a single
// webhook deployment**.
func (c *Solver) Name() string {
	return "hetzner"
}

// Initialize will be called when the webhook first starts. This method can be used to
// instantiate the webhook, i.e. initialising connections or warming up caches.
// Typically, the kubeClientConfig parameter is used to build a Kubernetes client that
// can be used to fetch resources from the Kubernetes API, e.g. Secret resources
// containing credentials used to authenticate with DNS provider accounts.
// The stopCh can be used to handle early termination of the webhook, in cases where a
// SIGTERM or similar signal is sent to the webhook process.
func (c *Solver) Initialize(kubeClientConfig *rest.Config, _ <-chan struct{}) error {
	cl, err := kubernetes.NewForConfig(kubeClientConfig)
	if err != nil {
		return err
	}

	c.hClientBuilder = NewHClientBuilder(cl, c.registry)

	return nil
}

// Present is responsible for actually presenting the DNS record with the DNS provider.
// This method should tolerate being called multiple times with the same value.
// cert-manager itself will later perform a self check to ensure that the solver has
// correctly configured the DNS provider.
func (c *Solver) Present(ch *v1alpha1.ChallengeRequest) error {
	ctx := context.Background()

	cfg, err := LoadConfig(ch.Config)
	if err != nil {
		return err
	}

	hClient, err := c.hClientBuilder(ctx, ch.ResourceNamespace, cfg)
	if err != nil {
		return err
	}

	zoneRRSet, err := BuildZoneRRSet(ch)
	if err != nil {
		return fmt.Errorf("error building zone and zone rrset: %w", err)
	}

	c.logger.Info(
		"creating DNS TXT record",
		"zone-name", zoneRRSet.Zone.Name,
		"zone-rrset-name", zoneRRSet.Name,
	)

	action, _, err := hClient.Zone.AddRRSetRecords(ctx,
		zoneRRSet,
		hcloud.ZoneRRSetAddRecordsOpts{
			Records: []hcloud.ZoneRRSetRecord{{Value: zoneutil.FormatTXTRecord(ch.Key)}},
			TTL:     hcloud.Ptr(TTL),
		},
	)
	if err != nil {
		errMsg := "failed to request rrset record addition"
		c.logger.Error(errMsg, "err", err)
		return fmt.Errorf("%s: %w", errMsg, hcloud.StabilizeError(err))
	}

	if err := hClient.Action.WaitFor(ctx, action); err != nil {
		errMsg := "failed to add rrset record"
		c.logger.Error(errMsg, "err", err)
		return fmt.Errorf("%s: %w", errMsg, hcloud.StabilizeError(err))
	}

	return nil
}

// CleanUp should delete the relevant TXT record from the DNS provider console.
// If multiple TXT records exist with the same record name (e.g.
// _acme-challenge.example.com) then **only** the record with the same `key` value
// provided on the ChallengeRequest should be cleaned up. This is in order to facilitate
// multiple DNS validations for the same domain concurrently.
func (c *Solver) CleanUp(ch *v1alpha1.ChallengeRequest) error {
	ctx := context.Background()

	cfg, err := LoadConfig(ch.Config)
	if err != nil {
		return err
	}

	hClient, err := c.hClientBuilder(ctx, ch.ResourceNamespace, cfg)
	if err != nil {
		return err
	}

	zoneRRSet, err := BuildZoneRRSet(ch)
	if err != nil {
		return fmt.Errorf("error building zone and zone rrset: %w", err)
	}

	c.logger.Info(
		"removing DNS TXT record",
		"zone-name", zoneRRSet.Zone.Name,
		"zone-rrset-name", zoneRRSet.Name,
	)

	action, _, err := hClient.Zone.RemoveRRSetRecords(ctx,
		zoneRRSet,
		hcloud.ZoneRRSetRemoveRecordsOpts{
			Records: []hcloud.ZoneRRSetRecord{{Value: fmt.Sprintf("%q", ch.Key)}},
		},
	)
	if err != nil {
		if hcloud.IsError(err, hcloud.ErrorCodeNotFound) {
			c.logger.Info(
				"zone rrset has already been deleted",
				"zone-name", zoneRRSet.Zone.Name,
				"zone-rrset-name", zoneRRSet.Name,
			)
			return nil
		}
		errMsg := "failed to request rrset record deletion"
		c.logger.Error(errMsg, "err", err)
		return fmt.Errorf("%s: %w", errMsg, hcloud.StabilizeError(err))
	}

	if err := hClient.Action.WaitFor(ctx, action); err != nil {
		errMsg := "failed to delete rrset record"
		c.logger.Error(errMsg, "err", err)
		return fmt.Errorf("%s: %w", errMsg, hcloud.StabilizeError(err))
	}

	return nil
}
