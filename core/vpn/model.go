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

type Destination interface {
	Hash() string
	VPNID() uuid.UUID
	SourceName() string
	Targets() []DestinationTarget
}

type DestinationTarget struct {
	Host  string
	IP    string
	Port  int
	HTTPS bool
}
