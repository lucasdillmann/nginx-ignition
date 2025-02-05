package integration

import "dillmann.com.br/nginx-ignition/api/common/dynamic_field"

type integrationConfigurationRequest struct {
	Enabled    *bool           `json:"enabled" validation:"required"`
	Parameters *map[string]any `json:"parameters" validation:"required"`
}

type integrationConfigurationResponse struct {
	ID                  string                                `json:"id"`
	Name                string                                `json:"name"`
	Description         string                                `json:"description"`
	Enabled             bool                                  `json:"enabled"`
	ConfigurationFields []*dynamic_field.DynamicFieldResponse `json:"configurationFields"`
	Parameters          map[string]any                        `json:"parameters"`
}

type integrationOptionResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type integrationResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
}
