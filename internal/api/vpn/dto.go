package vpn

import (
	"github.com/google/uuid"

	"nginx-ignition/internal/api/common/dynamicfield"
	"nginx-ignition/internal/core/common/i18n"
	"nginx-ignition/internal/core/vpn"
)

type vpnRequest struct {
	Parameters map[string]any `json:"parameters"`
	Name       string         `json:"name"`
	Driver     string         `json:"driver"`
	Enabled    bool           `json:"enabled"`
}

type vpnResponse struct {
	Parameters               map[string]any         `json:"parameters"`
	Name                     string                 `json:"name"`
	Driver                   string                 `json:"driver"`
	DriverEndpointSSLSupport vpn.EndpointSSLSupport `json:"driverEndpointSslSupport"`
	ID                       uuid.UUID              `json:"id"`
	Enabled                  bool                   `json:"enabled"`
}

type vpnDriverResponse struct {
	Name                  *i18n.Message           `json:"name"`
	ID                    string                  `json:"id"`
	EndpointSSLSupport    vpn.EndpointSSLSupport  `json:"endpointSslSupport"`
	ImportantInstructions []*i18n.Message         `json:"importantInstructions"`
	ConfigurationFields   []dynamicfield.Response `json:"configurationFields"`
}
