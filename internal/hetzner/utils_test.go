package hetzner

import (
	"testing"

	"github.com/cert-manager/cert-manager/pkg/acme/webhook/apis/acme/v1alpha1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestBuildZoneRRSet(t *testing.T) {
	tests := []struct {
		name          string
		ch            *v1alpha1.ChallengeRequest
		wantZoneRRSet *hcloud.ZoneRRSet
		wantErr       error
	}{
		{
			name: "success",
			ch:   &v1alpha1.ChallengeRequest{ResolvedZone: "example.com.", ResolvedFQDN: "_acme-challenge.example.com."},
			wantZoneRRSet: &hcloud.ZoneRRSet{
				Zone: &hcloud.Zone{Name: "example.com"},
				Name: "_acme-challenge",
				Type: hcloud.ZoneRRSetTypeTXT,
			},
			wantErr: nil,
		},

		{
			name: "success punycode",
			ch:   &v1alpha1.ChallengeRequest{ResolvedZone: "münchen.com.", ResolvedFQDN: "_acme-challenge.münchen.com."},
			wantZoneRRSet: &hcloud.ZoneRRSet{
				Zone: &hcloud.Zone{Name: "xn--mnchen-3ya.com"},
				Name: "_acme-challenge",
				Type: hcloud.ZoneRRSetTypeTXT,
			},
			wantErr: nil,
		},
		{
			name: "success with inconsistent dots",
			ch:   &v1alpha1.ChallengeRequest{ResolvedZone: "example.com", ResolvedFQDN: "_acme-challenge.example.com."},
			wantZoneRRSet: &hcloud.ZoneRRSet{
				Zone: &hcloud.Zone{Name: "example.com"},
				Name: "_acme-challenge",
				Type: hcloud.ZoneRRSetTypeTXT,
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			zoneRRSet, err := BuildZoneRRSet(tt.ch)
			if tt.wantErr != nil {
				require.EqualError(t, err, tt.wantErr.Error())
			}

			assert.Equal(t, tt.wantZoneRRSet, zoneRRSet)
		})
	}
}
