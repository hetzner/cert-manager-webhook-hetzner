package hetzner

import (
	"fmt"
	"strings"

	"github.com/cert-manager/cert-manager/pkg/acme/webhook/apis/acme/v1alpha1"
	"github.com/cert-manager/cert-manager/pkg/issuer/acme/dns/util"
	"golang.org/x/net/idna"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func BuildZoneRRSet(
	ch *v1alpha1.ChallengeRequest,
) (*hcloud.ZoneRRSet, error) {
	resolvedZone, err := idna.ToASCII(util.UnFqdn(ch.ResolvedZone))
	if err != nil {
		return nil, fmt.Errorf("error converting ResolvedZone to ASCII: %w", err)
	}

	resolvedFQDN, err := idna.ToASCII(util.UnFqdn(ch.ResolvedFQDN))
	if err != nil {
		return nil, fmt.Errorf("error converting ResolvedFQDN to ASCII: %w", err)
	}

	zoneRRSetName := strings.TrimSuffix(resolvedFQDN, "."+resolvedZone)

	zone := &hcloud.Zone{
		Name: resolvedZone,
	}

	return &hcloud.ZoneRRSet{
			Zone: zone,
			Name: zoneRRSetName,
			Type: hcloud.ZoneRRSetTypeTXT,
		},
		nil
}
