package integration

import (
	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/common/pagination"
	"dillmann.com.br/nginx-ignition/core/integration"
)

func newIntegrationRequest() integrationRequest {
	return integrationRequest{
		Name:    "Test Integration",
		Driver:  "test-driver",
		Enabled: true,
		Parameters: map[string]any{
			"key": "value",
		},
	}
}

func newIntegration() *integration.Integration {
	return &integration.Integration{
		ID:      uuid.New(),
		Name:    "Test Integration",
		Driver:  "test-driver",
		Enabled: true,
		Parameters: map[string]any{
			"key": "value",
		},
	}
}

func newIntegrationPage() *pagination.Page[integration.Integration] {
	return pagination.Of([]integration.Integration{
		*newIntegration(),
	})
}
