package integration

import (
	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/api/common/dynamicfield"
	"dillmann.com.br/nginx-ignition/core/integration"
)

func toDto(data *integration.Integration) *integrationResponse {
	return &integrationResponse{
		ID:         data.ID,
		Driver:     data.Driver,
		Name:       data.Name,
		Enabled:    data.Enabled,
		Parameters: data.Parameters,
	}
}

func toDomain(data *integrationRequest, id uuid.UUID) *integration.Integration {
	return &integration.Integration{
		ID:         id,
		Driver:     data.Driver,
		Name:       data.Name,
		Enabled:    data.Enabled,
		Parameters: data.Parameters,
	}
}

func fromDto(id uuid.UUID, data *integrationRequest) *integration.Integration {
	return &integration.Integration{
		ID:         id,
		Driver:     data.Driver,
		Name:       data.Name,
		Enabled:    data.Enabled,
		Parameters: data.Parameters,
	}
}

func toOptionDto(option *integration.DriverOption) *integrationOptionResponse {
	return &integrationOptionResponse{
		ID:        option.ID,
		Name:      option.Name,
		Port:      option.Port,
		Qualifier: option.Qualifier,
		Protocol:  string(option.Protocol),
	}
}

func toAvailableDriverDto(data *integration.AvailableDriver) *integrationDriverResponse {
	return &integrationDriverResponse{
		ID:                  data.ID,
		Name:                data.Name,
		Description:         data.Description,
		ConfigurationFields: dynamicfield.ToResponse(data.ConfigurationFields),
	}
}
