package integration

import (
	"dillmann.com.br/nginx-ignition/api/common/dynamic_field"
	"dillmann.com.br/nginx-ignition/core/integration"
)

func toDto(integration *integration.ListOutput) *integrationResponse {
	return &integrationResponse{
		ID:          integration.ID,
		Name:        integration.Name,
		Description: integration.Description,
		Enabled:     integration.Enabled,
	}
}

func toOptionDto(option *integration.AdapterOption) *integrationOptionResponse {
	return &integrationOptionResponse{
		ID:       option.ID,
		Name:     option.Name,
		Port:     option.Port,
		Protocol: string(option.Protocol),
	}
}

func toConfigurationDto(configuration *integration.GetByIdOutput) *integrationConfigurationResponse {
	return &integrationConfigurationResponse{
		ID:                  configuration.ID,
		Name:                configuration.Name,
		Description:         configuration.Description,
		Enabled:             configuration.Enabled,
		ConfigurationFields: dynamic_field.ToResponse(configuration.ConfigurationFields),
		Parameters:          configuration.Parameters,
	}
}
