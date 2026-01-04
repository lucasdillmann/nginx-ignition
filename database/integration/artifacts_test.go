package integration

import (
	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/integration"
)

func newIntegration() *integration.Integration {
	return &integration.Integration{
		ID:      uuid.New(),
		Name:    uuid.NewString(),
		Driver:  "DOCKER",
		Enabled: true,
		Parameters: map[string]any{
			"key": "value",
		},
	}
}
