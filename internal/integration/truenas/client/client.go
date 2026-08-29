package client

import (
	"github.com/lucasdillmann/nginx-ignition/internal/core/common/configuration"
	"github.com/lucasdillmann/nginx-ignition/internal/integration/truenas/fields"
)

type Client interface {
	GetAvailableApps() ([]AvailableAppDTO, error)
}

func For(cfg *configuration.Configuration, parameters map[string]any) (Client, error) {
	if err := initCache(cfg); err != nil {
		return nil, err
	}

	baseURL := parameters[fields.URLFieldID].(string)
	username := parameters[fields.UsernameFieldID].(string)
	password := parameters[fields.PasswordFieldID].(string)

	if useLegacyAPI(parameters) {
		return newRestClient(baseURL, username, password), nil
	}

	return newWebSocketClient(baseURL, username, password), nil
}

func useLegacyAPI(parameters map[string]any) bool {
	rawValue, found := parameters[fields.LegacyAPIFieldID]
	if !found {
		return true
	}

	parsedValue, parsed := rawValue.(bool)
	return parsed && parsedValue
}
