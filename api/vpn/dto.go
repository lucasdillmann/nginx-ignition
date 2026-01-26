package vpn

import (
	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/api/common/dynamicfield"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
)

type vpnRequest struct {
	Parameters map[string]any `json:"parameters"`
	Name       string         `json:"name"`
	Driver     string         `json:"driver"`
	Enabled    bool           `json:"enabled"`
}

type vpnResponse struct {
	Parameters map[string]any `json:"parameters"`
	Name       string         `json:"name"`
	Driver     string         `json:"driver"`
	ID         uuid.UUID      `json:"id"`
	Enabled    bool           `json:"enabled"`
}

type vpnDriverResponse struct {
	ID                    string                  `json:"id"`
	Name                  *i18n.Message           `json:"name"`
	ImportantInstructions []*i18n.Message         `json:"importantInstructions"`
	ConfigurationFields   []dynamicfield.Response `json:"configurationFields"`
}
