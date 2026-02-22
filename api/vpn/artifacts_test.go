package vpn

import (
	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/common/i18n"
	"dillmann.com.br/nginx-ignition/core/common/pagination"
	"dillmann.com.br/nginx-ignition/core/vpn"
)

func newVPN() *vpn.VPN {
	return &vpn.VPN{
		ID:      uuid.New(),
		Name:    "Test VPN",
		Driver:  "test-driver",
		Enabled: true,
		Parameters: map[string]any{
			"key": "value",
		},
	}
}

func newVPNAvailableDriver() *vpn.AvailableDriver {
	return &vpn.AvailableDriver{
		ID:                 "test-driver",
		Name:               i18n.Static("Test Driver"),
		EndpointSSLSupport: vpn.DriverManagedEndpointSSLSupport,
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

func newVPNPage() *pagination.Page[vpn.VPN] {
	return pagination.Of([]vpn.VPN{
		*newVPN(),
	})
}
