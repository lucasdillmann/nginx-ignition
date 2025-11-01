package vpn

import (
	"github.com/google/uuid"
)

type VPN struct {
	ID         uuid.UUID
	Driver     string
	Name       string
	Enabled    bool
	Parameters map[string]any
}
