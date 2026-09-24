package vpn

import (
	"github.com/google/uuid"

	"github.com/lucasdillmann/nginx-ignition/internal/business/core/i18n"
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/pagination"
	vpn2 "github.com/lucasdillmann/nginx-ignition/internal/business/domain/vpn"
)

func newVPN() *vpn2.VPN {
	return &vpn2.VPN{
		ID:      uuid.New(),
		Name:    "Test VPN",
		Driver:  "test-driver",
		Enabled: true,
		Parameters: map[string]any{
			"key": "value",
		},
	}
}

func newVPNAvailableDriver() *vpn2.AvailableDriver {
	return &vpn2.AvailableDriver{
		ID:                 "test-driver",
		Name:               i18n.Static("Test Driver"),
		EndpointSSLSupport: vpn2.DriverManagedEndpointSSLSupport,
	}
}

func newVPNRequest() vpnRequest {
	return vpnRequest{
		Name:    "Test VPN",
		Driver:  "test-driver",
		Enabled: true,
		Parameters: map[string]any{
			"key": "value",
		},
	}
}

func newVPNPage() *pagination.Page[vpn2.VPN] {
	return pagination.Of([]vpn2.VPN{
		*newVPN(),
	})
}
