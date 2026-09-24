package vpn

import (
	"github.com/google/uuid"

	"github.com/lucasdillmann/nginx-ignition/internal/api/core/dynamicfield"

	vpn2 "github.com/lucasdillmann/nginx-ignition/internal/business/domain/vpn"
)

func toDTO(data *vpn2.VPN, driver *vpn2.AvailableDriver) *vpnResponse {
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

func toDomain(data *vpnRequest, id uuid.UUID) *vpn2.VPN {
	return &vpn2.VPN{
		ID:         id,
		Driver:     data.Driver,
		Name:       data.Name,
		Enabled:    data.Enabled,
		Parameters: data.Parameters,
	}
}

func fromDTO(id uuid.UUID, data *vpnRequest) *vpn2.VPN {
	return &vpn2.VPN{
		ID:         id,
		Driver:     data.Driver,
		Name:       data.Name,
		Enabled:    data.Enabled,
		Parameters: data.Parameters,
	}
}

func toAvailableDriverDTO(data *vpn2.AvailableDriver) vpnDriverResponse {
	return vpnDriverResponse{
		ID:                    data.ID,
		Name:                  data.Name,
		ImportantInstructions: data.ImportantInstructions,
		ConfigurationFields:   dynamicfield.ToResponse(data.ConfigurationFields),
		EndpointSSLSupport:    data.EndpointSSLSupport,
	}
}
