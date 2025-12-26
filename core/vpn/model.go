package vpn

import (
	"github.com/google/uuid"
)

type VPN struct {
	Parameters map[string]any
	Driver     string
	Name       string
	ID         uuid.UUID
	Enabled    bool
}

type Endpoint interface {
	Hash() string
	VPNID() uuid.UUID
	SourceName() string
	Targets() []EndpointTarget
}

type EndpointTarget struct {
	Host  string
	IP    string
	Port  int
	HTTPS bool
}
