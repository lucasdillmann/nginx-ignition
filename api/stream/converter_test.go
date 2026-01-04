package stream

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"dillmann.com.br/nginx-ignition/core/stream"
)

func Test_toDTO(t *testing.T) {
	t.Run("converts domain object to DTO", func(t *testing.T) {
		subject := newStream()
		result := toDTO(subject)

		assert.NotNil(t, result)
		assert.Equal(t, subject.ID, *result.ID)
		assert.True(t, *result.Enabled)
		assert.Equal(t, subject.Name, *result.Name)
		assert.Equal(t, string(subject.Type), *result.Type)
		assert.Equal(t, subject.Binding.Address, *result.Binding.Address)
		assert.Equal(t, *subject.Binding.Port, *result.Binding.Port)
		assert.Equal(
			t,
			subject.DefaultBackend.Address.Address,
			*result.DefaultBackend.Target.Address,
		)
	})

	t.Run("returns nil when input is nil", func(t *testing.T) {
		result := toDTO(nil)
		assert.Nil(t, result)
	})
}

func Test_toDomain(t *testing.T) {
	t.Run("converts DTO to domain object", func(t *testing.T) {
		payload := newStreamRequest()
		result := toDomain(&payload)

		assert.NotNil(t, result)
		assert.True(t, result.Enabled)
		assert.Equal(t, *payload.Name, result.Name)
		assert.Equal(t, stream.SimpleType, result.Type)
		assert.Equal(t, *payload.Binding.Address, result.Binding.Address)
		assert.Equal(t, *payload.Binding.Port, *result.Binding.Port)
		assert.Equal(
			t,
			*payload.DefaultBackend.Target.Address,
			result.DefaultBackend.Address.Address,
		)
	})

	t.Run("returns nil when input is nil", func(t *testing.T) {
		result := toDomain(nil)
		assert.Nil(t, result)
	})
}
