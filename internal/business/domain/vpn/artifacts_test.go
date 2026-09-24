package vpn

import (
	"github.com/google/uuid"
)

func newVPN() *VPN {
	return &VPN{
		ID:         uuid.New(),
		Name:       "test",
		Driver:     "tailscale",
		Enabled:    true,
		Parameters: map[string]any{},
	}
}
