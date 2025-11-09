package vpn

import (
	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/api/common/dynamic_field"
	"dillmann.com.br/nginx-ignition/core/vpn"
)

func toDto(data *vpn.VPN) *vpnResponse {
	return &vpnResponse{
		ID:         data.ID,
		Driver:     data.Driver,
		Name:       data.Name,
		Enabled:    data.Enabled,
		Parameters: data.Parameters,
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

func fromDto(id uuid.UUID, data *vpnRequest) *vpn.VPN {
	return &vpn.VPN{
		ID:         id,
		Driver:     data.Driver,
		Name:       data.Name,
		Enabled:    data.Enabled,
		Parameters: data.Parameters,
	}
}

func toAvailableDriverDto(data *vpn.AvailableDriver) *vpnDriverResponse {
	return &vpnDriverResponse{
		ID:                    data.ID,
		Name:                  data.Name,
		ImportantInstructions: data.ImportantInstructions,
		ConfigurationFields:   dynamic_field.ToResponse(data.ConfigurationFields),
	}
}
