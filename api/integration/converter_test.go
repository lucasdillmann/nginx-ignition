package integration

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"dillmann.com.br/nginx-ignition/core/integration"
)

func Test_ToDTO(t *testing.T) {
	t.Run("converts domain object to DTO", func(t *testing.T) {
		id := uuid.New()
		input := &integration.Integration{
			ID:         id,
			Name:       "integration-1",
			Driver:     "docker",
			Enabled:    true,
			Parameters: map[string]any{"key": "value"},
		}

		result := toDTO(input)

		assert.NotNil(t, result)
		assert.Equal(t, id, result.ID)
		assert.Equal(t, input.Name, result.Name)
		assert.Equal(t, input.Driver, result.Driver)
		assert.True(t, result.Enabled)
		assert.Equal(t, input.Parameters, result.Parameters)
	})

	t.Run("returns nil when input is nil", func(t *testing.T) {
		result := toDTO(nil)
		assert.Nil(t, result)
	})
}

func Test_ToDomain(t *testing.T) {
	t.Run("converts DTO to domain object", func(t *testing.T) {
		id := uuid.New()
		input := &integrationRequest{
			Name:       "integration-1",
			Driver:     "docker",
			Enabled:    true,
			Parameters: map[string]any{"key": "value"},
		}

		result := toDomain(input, id)

		assert.NotNil(t, result)
		assert.Equal(t, id, result.ID)
		assert.Equal(t, input.Name, result.Name)
		assert.Equal(t, input.Driver, result.Driver)
		assert.True(t, result.Enabled)
		assert.Equal(t, input.Parameters, result.Parameters)
	})
}
