package integration

import (
	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/api/common/dynamicfield"
)

type integrationRequest struct {
	Name       string         `json:"name"`
	Driver     string         `json:"driver"`
	Enabled    bool           `json:"enabled"`
	Parameters map[string]any `json:"parameters"`
}

type integrationResponse struct {
	ID         uuid.UUID      `json:"id"`
	Name       string         `json:"name"`
	Driver     string         `json:"driver"`
	Enabled    bool           `json:"enabled"`
	Parameters map[string]any `json:"parameters"`
}

type integrationOptionResponse struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Port      int     `json:"port"`
	Qualifier *string `json:"qualifier"`
	Protocol  string  `json:"protocol"`
}

type integrationDriverResponse struct {
	ID                  string                              `json:"id"`
	Name                string                              `json:"name"`
	Description         string                              `json:"description"`
	ConfigurationFields []dynamicfield.DynamicFieldResponse `json:"configurationFields"`
}
