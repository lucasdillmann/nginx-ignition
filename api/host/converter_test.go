package host

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"dillmann.com.br/nginx-ignition/core/binding"
	"dillmann.com.br/nginx-ignition/core/common/ptr"
	"dillmann.com.br/nginx-ignition/core/host"
	"dillmann.com.br/nginx-ignition/core/settings"
)

func Test_ToDTO(t *testing.T) {
	t.Run("converts domain object to DTO", func(t *testing.T) {
		id := uuid.New()
		input := &host.Host{
			ID:                id,
			Enabled:           true,
			DefaultServer:     true,
			UseGlobalBindings: true,
			DomainNames:       []string{"example.com"},
			FeatureSet: host.FeatureSet{
				WebsocketSupport: true,
			},
		}

		globalSettings := &settings.Settings{
			GlobalBindings: []binding.Binding{
				{Type: binding.HTTPBindingType, Port: 80},
			},
		}

		result := toDTO(input, globalSettings)

		assert.NotNil(t, result)
		assert.Equal(t, id, *result.ID)
		assert.True(t, *result.Enabled)
		assert.True(t, *result.DefaultServer)
		assert.True(t, *result.UseGlobalBindings)
		assert.Equal(t, input.DomainNames, result.DomainNames)
		assert.Len(t, result.GlobalBindings, 1)
		assert.True(t, *result.FeatureSet.WebsocketsSupport)
	})

	t.Run("returns nil when input is nil", func(t *testing.T) {
		result := toDTO(nil, nil)
		assert.Nil(t, result)
	})
}

func Test_ToDomain(t *testing.T) {
	t.Run("converts DTO to domain object", func(t *testing.T) {
		input := &hostRequestDTO{
			Enabled:           ptr.Of(true),
			DefaultServer:     ptr.Of(true),
			UseGlobalBindings: ptr.Of(false),
			DomainNames:       []string{"example.com"},
			FeatureSet: &featureSetDTO{
				WebsocketsSupport: ptr.Of(true),
			},
		}

		result := toDomain(input)

		assert.NotNil(t, result)
		assert.True(t, result.Enabled)
		assert.True(t, result.DefaultServer)
		assert.False(t, result.UseGlobalBindings)
		assert.Equal(t, input.DomainNames, result.DomainNames)
		assert.True(t, result.FeatureSet.WebsocketSupport)
	})

	t.Run("returns nil when input is nil", func(t *testing.T) {
		result := toDomain(nil)
		assert.Nil(t, result)
	})
}
