package vpn

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"dillmann.com.br/nginx-ignition/core/vpn"
)

func Test_Converter(t *testing.T) {
	t.Run("toDTO", func(t *testing.T) {
		t.Run("converts domain object to DTO", func(t *testing.T) {
			id := uuid.New()
			input := &vpn.VPN{
				ID:         id,
				Name:       "vpn-1",
				Driver:     "openvpn",
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
	})

	t.Run("toDomain", func(t *testing.T) {
		t.Run("converts DTO to domain object", func(t *testing.T) {
			id := uuid.New()
			input := &vpnRequest{
				Name:       "vpn-1",
				Driver:     "openvpn",
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
	})
}
