package hetzner

import (
	"encoding/json"
	"fmt"

	extapi "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

// Config is a structure that is used to decode into when solving a DNS01 challenge.
// This information is provided by cert-manager, and may be a reference to additional
// configuration that's needed to solve the challenge for this particular certificate or
// issuer. This typically includes references to Secret resources containing DNS
// provider credentials, in cases where a 'multi-tenant' DNS solver is being created.
type Config struct {
	HetznerTokenSecret   SecretKeyRef `json:"tokenSecretKeyRef"`
	HetznerTokenFilePath string       `json:"tokenFilePath"`
	HCloudEndpoint       string       `json:"hcloudEndpoint"`

	// HTTPTimeout overrides the default per-request timeout used for calls to the
	// Hetzner Cloud API, e.g. "45s". Must be a positive duration parseable by
	// [time.ParseDuration]. Falls back to [DefaultHTTPTimeout] when empty.
	HTTPTimeout string `json:"httpTimeout"`
}

type SecretKeyRef struct {
	Name string `json:"name"`
	Key  string `json:"key"`
}

// LoadConfig is a small helper function that decodes JSON configuration into the typed
// config struct.
func LoadConfig(cfgJSON *extapi.JSON) (Config, error) {
	cfg := Config{}
	// handle the 'base case' where no configuration has been provided
	if cfgJSON == nil {
		return cfg, nil
	}
	if err := json.Unmarshal(cfgJSON.Raw, &cfg); err != nil {
		return cfg, fmt.Errorf("error decoding solver config: %w", err)
	}

	return cfg, nil
}
