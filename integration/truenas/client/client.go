package client

import (
	"dillmann.com.br/nginx-ignition/core/common/configuration"
	"dillmann.com.br/nginx-ignition/integration/truenas/fields"
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
		return false
	}

	parsedValue, parsed := rawValue.(bool)
	return parsed && parsedValue
}
