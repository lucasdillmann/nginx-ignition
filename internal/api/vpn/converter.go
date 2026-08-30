package vpn

import (
	"github.com/google/uuid"

	"github.com/lucasdillmann/nginx-ignition/internal/api/common/dynamicfield"
	"github.com/lucasdillmann/nginx-ignition/internal/core/vpn"
)

func toDTO(data *vpn.VPN, driver *vpn.AvailableDriver) *vpnResponse {
	if data == nil {
		return nil
	}

	return &vpnResponse{
		ID:                       data.ID,
		Driver:                   data.Driver,
		DriverEndpointSSLSupport: driver.EndpointSSLSupport,
		Name:                     data.Name,
		Enabled:                  data.Enabled,
		Parameters:               data.Parameters,
	}
}

func toDomain(data *vpnRequest, id uuid.UUID) *vpn.VPN {
	return &vpn.VPN{
		ID:         id,
		Driver:     data.Driver,
		Name:       data.Name,
		Enabled:    data.Enabled,
		Parameters: data.Parameters,
	}
}

func fromDTO(id uuid.UUID, data *vpnRequest) *vpn.VPN {
	return &vpn.VPN{
		ID:         id,
		Driver:     data.Driver,
		Name:       data.Name,
		Enabled:    data.Enabled,
		Parameters: data.Parameters,
	}
}

func toAvailableDriverDTO(data *vpn.AvailableDriver) vpnDriverResponse {
	return vpnDriverResponse{
		ID:                    data.ID,
		Name:                  data.Name,
		ImportantInstructions: data.ImportantInstructions,
		ConfigurationFields:   dynamicfield.ToResponse(data.ConfigurationFields),
		EndpointSSLSupport:    data.EndpointSSLSupport,
	}
}
