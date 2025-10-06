//go:build e2e

package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"testing"
	"time"

	acmetest "github.com/cert-manager/cert-manager/test/acme"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/require"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	"github.com/hetzner/cert-manager-webhook-hetzner/internal/hetzner"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/hetznercloud/hcloud-go/v2/hcloud/exp/kit/envutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud/exp/kit/randutil"
)

func init() {
	// Do no panic when the controller logger has not been configured.
	ctrl.SetLogger(zap.New(zap.WriteTo(os.Stderr)))
}

func TestRunsSuite(t *testing.T) {
	ctx := t.Context()

	hetznerToken, err := envutil.LookupEnvWithFile("HETZNER_TOKEN")
	require.NoError(t, err)
	hcloudEndpoint, err := envutil.LookupEnvWithFile("HCLOUD_ENDPOINT")
	require.NoError(t, err)

	clientOpts := []hcloud.ClientOption{
		hcloud.WithToken(hetznerToken),
	}
	if hcloudEndpoint != "" {
		clientOpts = append(clientOpts, hcloud.WithEndpoint(hcloudEndpoint))
	}

	cl := hcloud.NewClient(clientOpts...)

	// Create random run data
	runID := randutil.GenerateID()
	zoneName := fmt.Sprintf("example-%s.com", runID)
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

	// Create test zone in hcloud
	result, _, err := cl.Zone.Create(ctx, hcloud.ZoneCreateOpts{
		Name: zoneName,
		Mode: hcloud.ZoneModePrimary,
	})
	require.NoError(t, err)
	require.NoError(t, cl.Action.WaitFor(ctx, result.Action))

	backoffFn := hcloud.ExponentialBackoffWithOpts(
		hcloud.ExponentialBackoffOpts{
			Base:       time.Millisecond * 250,
			Cap:        time.Second * 60,
			Multiplier: 2.0,
		},
	)

	// Wait for the delegation to finish
	zone := result.Zone
	retries := 0
	for retries < 10 {
		if zone.AuthoritativeNameservers.DelegationStatus != hcloud.ZoneDelegationStatusUnknown {
			break
		}

		zone, _, err = cl.Zone.Get(ctx, zone.Name)
		require.NoError(t, err)

		t.Logf("waiting for zone delegation to finish (attempt %d)", retries)
		select {
		case <-ctx.Done():
		case <-time.After(backoffFn(retries)):
		}

		retries++
	}

	// Register cleanup test zone in hcloud
	t.Cleanup(func() {
		cleanCtx := context.Background()
		result, _, err := cl.Zone.Delete(cleanCtx, result.Zone)
		require.NoError(t, err)
		require.NoError(t, cl.Action.WaitFor(cleanCtx, result.Action))
	})

	// Fetch DNS server address
	zoneRefreshed, _, err := cl.Zone.Get(ctx, result.Zone.Name)
	require.NoError(t, err)
	require.NotEmpty(t, zoneRefreshed.AuthoritativeNameservers.Assigned)
	dnsServerAddr := net.JoinHostPort(
		zoneRefreshed.AuthoritativeNameservers.Assigned[0],
		"53",
	)

	// Run tests
	// Once https://github.com/cert-manager/cert-manager/pull/4835 is merged,
	// only run fixture.RunConformance(t) instead of fixture.RunBasic(t) and
	// fixture.RunExtended(t)
	fixture := acmetest.NewFixture(
		hetzner.New(logger, prometheus.DefaultRegisterer),
		acmetest.SetDNSName(zoneRefreshed.Name),
		acmetest.SetResolvedZone(fmt.Sprintf("%s.", zoneRefreshed.Name)),
		acmetest.SetAllowAmbientCredentials(false),
		acmetest.SetManifestPath("testdata/hetzner"),
		acmetest.SetDNSServer(dnsServerAddr),
		acmetest.SetUseAuthoritative(true),
		acmetest.SetStrict(true),
	)

	fixture.RunBasic(t)
	fixture.RunExtended(t)
}
