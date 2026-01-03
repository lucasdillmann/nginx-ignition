package healthcheck

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"dillmann.com.br/nginx-ignition/core/common/healthcheck"
)

func Test_Converter(t *testing.T) {
	t.Run("toDTO", func(t *testing.T) {
		t.Run("converts healthy status to DTO", func(t *testing.T) {
			status := &healthcheck.Status{
				Healthy: true,
				Details: []healthcheck.Detail{
					{ID: "db", Error: nil},
					{ID: "nginx", Error: nil},
				},
			}

			result := toDTO(status)

			assert.NotNil(t, result)
			assert.True(t, result.Healthy)
			assert.Len(t, result.Details, 2)
			assert.Equal(t, "db", result.Details[0].Component)
			assert.True(t, result.Details[0].Healthy)
			assert.Equal(t, "nginx", result.Details[1].Component)
			assert.True(t, result.Details[1].Healthy)
		})

		t.Run("converts unhealthy status to DTO", func(t *testing.T) {
			status := &healthcheck.Status{
				Healthy: false,
				Details: []healthcheck.Detail{
					{ID: "db", Error: errors.New("db error")},
					{ID: "nginx", Error: nil},
				},
			}

			result := toDTO(status)

			assert.NotNil(t, result)
			assert.False(t, result.Healthy)
			assert.Len(t, result.Details, 2)
			assert.Equal(t, "db", result.Details[0].Component)
			assert.False(t, result.Details[0].Healthy)
			assert.Equal(t, "nginx", result.Details[1].Component)
			assert.True(t, result.Details[1].Healthy)
		})
	})
}
