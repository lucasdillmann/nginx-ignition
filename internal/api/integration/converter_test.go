package integration

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_toDTO(t *testing.T) {
	t.Run("converts domain object to DTO", func(t *testing.T) {
		input := newIntegration()
		result := toDTO(input)

		assert.NotNil(t, result)
		assert.Equal(t, input.ID, result.ID)
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

func Test_toDomain(t *testing.T) {
	t.Run("converts DTO to domain object", func(t *testing.T) {
		id := uuid.New()
		input := newIntegrationRequest()
		result := toDomain(&input, id)

		assert.NotNil(t, result)
		assert.Equal(t, id, result.ID)
		assert.Equal(t, input.Name, result.Name)
		assert.Equal(t, input.Driver, result.Driver)
		assert.True(t, result.Enabled)
		assert.Equal(t, input.Parameters, result.Parameters)
	})
}
