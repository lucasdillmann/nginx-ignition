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
	HTTPS EndpointHTTPS
	Port  int
}

type EndpointHTTPS struct {
	PrivateKey         string
	PublicKey          string
	CertificationChain []string
	Enabled            bool
}

type EndpointSSLSupport string

const (
	DriverManagedEndpointSSLSupport   EndpointSSLSupport = "DRIVER_MANAGED"
	ProviderManagedEndpointSSLSupport EndpointSSLSupport = "PROVIDER_MANAGED"
)
