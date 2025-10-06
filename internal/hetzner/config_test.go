package hetzner

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	extapi "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name    string
		raw     string
		want    Config
		wantErr error
	}{
		{
			name: "valid",
			raw: `{
				"hcloudEndpoint": "https://changed.com/v2",
				"tokenSecretKeyRef": {
					"name": "hetzner",
					"key": "token"
				}
			}`,
			want: Config{
				HetznerTokenSecret: SecretKeyRef{
					Name: "hetzner",
					Key:  "token",
				},
				HCloudEndpoint: "https://changed.com/v2",
			},
		},
		{
			name: "empty config",
			raw:  "",
			want: Config{},
		},
		{
			name:    "broken json config",
			raw:     "{",
			want:    Config{},
			wantErr: errors.New("error decoding solver config: unexpected end of JSON input"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var cfg *extapi.JSON
			if tt.raw != "" {
				cfg = &extapi.JSON{Raw: []byte(tt.raw)}
			}

			got, err := LoadConfig(cfg)
			if tt.wantErr != nil {
				require.EqualError(t, tt.wantErr, err.Error())
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
