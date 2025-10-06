package hetzner

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"testing"

	"github.com/cert-manager/cert-manager/pkg/acme/webhook/apis/acme/v1alpha1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/hetznercloud/hcloud-go/v2/hcloud/exp/kit/randutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud/exp/mockutil"
)

func TestPresent(t *testing.T) {
	testCases := []struct {
		name       string
		zoneNameFn func(id string) string
	}{
		{
			name:       "success",
			zoneNameFn: func(id string) string { return fmt.Sprintf("example-%s.com", id) },
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			server, client := makeTestUtils(t)

			runID := randutil.GenerateID()
			zoneName := testCase.zoneNameFn(runID)
			logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

			key := "329ef10055b46b3cbc57"

			o := &Solver{
				logger:         logger,
				hClientBuilder: MockHClientBuilder(client),
			}

			ch := &v1alpha1.ChallengeRequest{
				UID:               "61616bdf-44f1-4795-916d-4d1d05e7d1ad",
				Action:            v1alpha1.ChallengeActionPresent,
				Type:              "dns-01",
				DNSName:           zoneName,
				ResolvedFQDN:      fmt.Sprintf("_acme-challenge.%s.", zoneName),
				ResolvedZone:      fmt.Sprintf("%s.", zoneName),
				Key:               key,
				ResourceNamespace: "default",
			}

			server.Expect([]mockutil.Request{
				{
					Method: "POST", Path: fmt.Sprintf("/zones/%s/rrsets/_acme-challenge/TXT/actions/add_records", zoneName),
					Want: func(t *testing.T, r *http.Request) {
						body, err := io.ReadAll(r.Body)
						require.NoError(t, err)
						assert.JSONEq(t,
							fmt.Sprintf(`{"records": [{ "value": %q }], "ttl": 300}`, fmt.Sprintf("%q", key)),
							string(body),
						)
					},
					Status: 201,
					JSONRaw: `{
						"action": {"id": 12, "status": "running"}
					}`,
				},
				{
					Method: "GET", Path: "/actions?id=12&page=1&sort=status&sort=id",
					Status: 200,
					JSONRaw: `{
						"actions": [{"id": 12, "status": "success"}]
					}`,
				},
			})

			err := o.Present(ch)
			require.NoError(t, err)
		})
	}
}

func TestCleanup(t *testing.T) {
	testCases := []struct {
		name       string
		zoneNameFn func(id string) string
	}{
		{
			name:       "success",
			zoneNameFn: func(id string) string { return fmt.Sprintf("example-%s.com", id) },
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			server, client := makeTestUtils(t)

			runID := randutil.GenerateID()
			zoneName := testCase.zoneNameFn(runID)
			logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

			key := fmt.Sprintf("329ef10055b46b3cbc57%s", runID)

			o := &Solver{
				logger:         logger,
				hClientBuilder: MockHClientBuilder(client),
			}

			ch := &v1alpha1.ChallengeRequest{
				UID:               "61616bdf-44f1-4795-916d-4d1d05e7d1ad",
				Action:            v1alpha1.ChallengeActionPresent,
				Type:              "dns-01",
				DNSName:           zoneName,
				ResolvedFQDN:      fmt.Sprintf("_acme-challenge.%s.", zoneName),
				ResolvedZone:      fmt.Sprintf("%s.", zoneName),
				Key:               key,
				ResourceNamespace: "default",
			}

			server.Expect([]mockutil.Request{
				{
					Method: "POST", Path: fmt.Sprintf("/zones/%s/rrsets/_acme-challenge/TXT/actions/remove_records", zoneName),
					Want: func(t *testing.T, r *http.Request) {
						body, err := io.ReadAll(r.Body)
						require.NoError(t, err)
						assert.JSONEq(t,
							fmt.Sprintf(`{"records": [{ "value": %q }]}`, fmt.Sprintf("%q", key)),
							string(body),
						)
					},
					Status: 201,
					JSONRaw: `{
						"action": {"id": 12, "status": "running"}
					}`,
				},
				{
					Method: "GET", Path: "/actions?id=12&page=1&sort=status&sort=id",
					Status: 200,
					JSONRaw: `{
						"actions": [{"id": 12, "status": "success"}]
					}`,
				},
			})

			err := o.CleanUp(ch)
			require.NoError(t, err)
		})
	}
}

func makeTestUtils(t *testing.T) (*mockutil.Server, *hcloud.Client) {
	server := mockutil.NewServer(t, nil)

	client := hcloud.NewClient(
		hcloud.WithEndpoint(server.URL),
		hcloud.WithRetryOpts(hcloud.RetryOpts{BackoffFunc: hcloud.ConstantBackoff(0), MaxRetries: 5}),
		hcloud.WithPollOpts(hcloud.PollOpts{BackoffFunc: hcloud.ConstantBackoff(0)}),
	)

	return server, client
}
