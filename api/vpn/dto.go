package vpn

import (
	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/api/common/dynamic_field"
)

type vpnRequest struct {
	Name       string         `json:"name" validate:"required"`
	Driver     string         `json:"driver" validate:"required"`
	Enabled    bool           `json:"enabled" validate:"required"`
	Parameters map[string]any `json:"parameters" validate:"required"`
}

type vpnResponse struct {
	ID         uuid.UUID      `json:"id"`
	Name       string         `json:"name"`
	Driver     string         `json:"driver"`
	Enabled    bool           `json:"enabled"`
	Parameters map[string]any `json:"parameters"`
}

type vpnDriverResponse struct {
	ID                  string                                `json:"id"`
	Name                string                                `json:"name"`
	ConfigurationFields []*dynamic_field.DynamicFieldResponse `json:"configurationFields"`
}
