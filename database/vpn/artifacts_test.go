package vpn

import (
	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/vpn"
)

func newVPN() *vpn.VPN {
	return &vpn.VPN{
		ID:      uuid.New(),
		Name:    uuid.NewString(),
		Driver:  "TAILSCALE",
		Enabled: true,
		Parameters: map[string]any{
			"test": true,
		},
	}
}
