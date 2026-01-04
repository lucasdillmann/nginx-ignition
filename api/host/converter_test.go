package host

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"dillmann.com.br/nginx-ignition/core/binding"
	"dillmann.com.br/nginx-ignition/core/settings"
)

func Test_toDTO(t *testing.T) {
	t.Run("converts domain object to DTO", func(t *testing.T) {
		input := newHost()
		globalSettings := &settings.Settings{
			GlobalBindings: []binding.Binding{
				{
					Type: binding.HTTPBindingType,
					Port: 80,
				},
			},
		}

		result := toDTO(input, globalSettings)

		assert.NotNil(t, result)
		assert.Equal(t, input.ID, *result.ID)
		assert.True(t, *result.Enabled)
		assert.False(t, *result.DefaultServer)
		assert.True(t, *result.UseGlobalBindings)
		assert.Equal(t, input.DomainNames, result.DomainNames)
		assert.Len(t, result.GlobalBindings, 1)
		assert.True(t, *result.FeatureSet.WebsocketsSupport)
		assert.Equal(t, "index.html", *result.Routes[0].Settings.IndexFile)
	})

	t.Run("returns nil when input is nil", func(t *testing.T) {
		result := toDTO(nil, nil)
		assert.Nil(t, result)
	})
}

func Test_toDomain(t *testing.T) {
	t.Run("converts DTO to domain object", func(t *testing.T) {
		input := newHostRequestDTO()
		result := toDomain(&input)

		assert.NotNil(t, result)
		assert.True(t, result.Enabled)
		assert.False(t, result.DefaultServer)
		assert.True(t, result.UseGlobalBindings)
		assert.Equal(t, input.DomainNames, result.DomainNames)
		assert.True(t, result.FeatureSet.WebsocketSupport)
		assert.Equal(t, "index.html", *result.Routes[0].Settings.IndexFile)
	})

	t.Run("returns nil when input is nil", func(t *testing.T) {
		result := toDomain(nil)
		assert.Nil(t, result)
	})
}
