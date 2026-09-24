package integration

import (
	"github.com/google/uuid"

	"github.com/lucasdillmann/nginx-ignition/internal/api/core/dynamicfield"

	integration2 "github.com/lucasdillmann/nginx-ignition/internal/business/domain/integration"
)

func toDTO(data *integration2.Integration) *integrationResponse {
	if data == nil {
		return nil
	}

	return &integrationResponse{
		ID:         data.ID,
		Driver:     data.Driver,
		Name:       data.Name,
		Enabled:    data.Enabled,
		Parameters: data.Parameters,
	}
}

func toDomain(data *integrationRequest, id uuid.UUID) *integration2.Integration {
	return &integration2.Integration{
		ID:         id,
		Driver:     data.Driver,
		Name:       data.Name,
		Enabled:    data.Enabled,
		Parameters: data.Parameters,
	}
}

func fromDTO(id uuid.UUID, data *integrationRequest) *integration2.Integration {
	return &integration2.Integration{
		ID:         id,
		Driver:     data.Driver,
		Name:       data.Name,
		Enabled:    data.Enabled,
		Parameters: data.Parameters,
	}
}

func toOptionDTO(option *integration2.DriverOption) *integrationOptionResponse {
	return &integrationOptionResponse{
		ID:        option.ID,
		Name:      option.Name,
		Port:      option.Port,
		Qualifier: option.Qualifier,
		Protocol:  string(option.Protocol),
	}
}

func toAvailableDriverDTO(data *integration2.AvailableDriver) integrationDriverResponse {
	return integrationDriverResponse{
		ID:                  data.ID,
		Name:                data.Name,
		Description:         data.Description,
		ConfigurationFields: dynamicfield.ToResponse(data.ConfigurationFields),
	}
}
