package vpn

import (
	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/api/common/dynamic_field"
)

type vpnRequest struct {
	Name       string         `json:"name"`
	Driver     string         `json:"driver"`
	Enabled    bool           `json:"enabled"`
	Parameters map[string]any `json:"parameters"`
}

type vpnResponse struct {
	ID         uuid.UUID      `json:"id"`
	Name       string         `json:"name"`
	Driver     string         `json:"driver"`
	Enabled    bool           `json:"enabled"`
	Parameters map[string]any `json:"parameters"`
}

type vpnDriverResponse struct {
	ID                    string                                `json:"id"`
	Name                  string                                `json:"name"`
	ImportantInstructions []string                              `json:"importantInstructions"`
	ConfigurationFields   []*dynamic_field.DynamicFieldResponse `json:"configurationFields"`
}
