package integration

import (
	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/api/common/dynamicfield"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
)

type integrationRequest struct {
	Parameters map[string]any `json:"parameters"`
	Name       string         `json:"name"`
	Driver     string         `json:"driver"`
	Enabled    bool           `json:"enabled"`
}

type integrationResponse struct {
	Parameters map[string]any `json:"parameters"`
	Name       string         `json:"name"`
	Driver     string         `json:"driver"`
	ID         uuid.UUID      `json:"id"`
	Enabled    bool           `json:"enabled"`
}

type integrationOptionResponse struct {
	Qualifier *string `json:"qualifier"`
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Protocol  string  `json:"protocol"`
	Port      int     `json:"port"`
}

type integrationDriverResponse struct {
	ID                  string                  `json:"id"`
	Name                *i18n.Message           `json:"name"`
	Description         *i18n.Message           `json:"description"`
	ConfigurationFields []dynamicfield.Response `json:"configurationFields"`
}
