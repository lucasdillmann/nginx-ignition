package integration

import (
	"github.com/google/uuid"
)

func newIntegration() *Integration {
	return &Integration{
		ID:         uuid.New(),
		Name:       "test",
		Driver:     "docker",
		Enabled:    true,
		Parameters: map[string]any{},
	}
}
