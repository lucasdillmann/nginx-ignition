package stream

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"dillmann.com.br/nginx-ignition/core/common/ptr"
	"dillmann.com.br/nginx-ignition/core/stream"
)

func Test_ToDTO(t *testing.T) {
	t.Run("converts domain object to DTO", func(t *testing.T) {
		id := uuid.New()
		input := &stream.Stream{
			ID:      id,
			Enabled: true,
			Name:    "stream-1",
			Type:    stream.SimpleType,
			FeatureSet: stream.FeatureSet{
				UseProxyProtocol: true,
			},
			Binding: stream.Address{
				Address:  "0.0.0.0",
				Port:     ptr.Of(80),
				Protocol: stream.TCPProtocol,
			},
			DefaultBackend: stream.Backend{
				Address: stream.Address{
					Address: "1.1.1.1",
					Port:    ptr.Of(80),
				},
				Weight: ptr.Of(5),
			},
		}

		result := toDTO(input)

		assert.NotNil(t, result)
		assert.Equal(t, id, *result.ID)
		assert.True(t, *result.Enabled)
		assert.Equal(t, input.Name, *result.Name)
		assert.Equal(t, string(input.Type), *result.Type)
		assert.True(t, *result.FeatureSet.UseProxyProtocol)
		assert.Equal(t, input.Binding.Address, *result.Binding.Address)
		assert.Equal(t, *input.Binding.Port, *result.Binding.Port)
		assert.Equal(t, input.DefaultBackend.Address.Address, *result.DefaultBackend.Target.Address)
		assert.Equal(t, *input.DefaultBackend.Weight, *result.DefaultBackend.Weight)
	})

	t.Run("returns nil when input is nil", func(t *testing.T) {
		result := toDTO(nil)
		assert.Nil(t, result)
	})
}

func Test_ToDomain(t *testing.T) {
	t.Run("converts DTO to domain object", func(t *testing.T) {
		input := &streamRequestDTO{
			Enabled: ptr.Of(true),
			Name:    ptr.Of("stream-1"),
			Type:    ptr.Of(string(stream.SimpleType)),
			FeatureSet: &featureSetDTO{
				UseProxyProtocol: ptr.Of(true),
			},
			Binding: &addressDTO{
				Address:  ptr.Of("0.0.0.0"),
				Port:     ptr.Of(80),
				Protocol: stream.TCPProtocol,
			},
			DefaultBackend: &backendDTO{
				Target: &addressDTO{
					Address: ptr.Of("1.1.1.1"),
					Port:    ptr.Of(80),
				},
				Weight: ptr.Of(5),
			},
		}

		result := toDomain(input)

		assert.NotNil(t, result)
		assert.True(t, result.Enabled)
		assert.Equal(t, *input.Name, result.Name)
		assert.Equal(t, stream.SimpleType, result.Type)
		assert.True(t, result.FeatureSet.UseProxyProtocol)
		assert.Equal(t, *input.Binding.Address, result.Binding.Address)
		assert.Equal(t, *input.Binding.Port, *result.Binding.Port)
		assert.Equal(t, *input.DefaultBackend.Target.Address, result.DefaultBackend.Address.Address)
		assert.Equal(t, *input.DefaultBackend.Weight, *result.DefaultBackend.Weight)
	})

	t.Run("returns nil when input is nil", func(t *testing.T) {
		result := toDomain(nil)
		assert.Nil(t, result)
	})
}
