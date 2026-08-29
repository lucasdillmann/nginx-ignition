package healthcheck

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_toDTO(t *testing.T) {
	t.Run("converts healthy status to DTO", func(t *testing.T) {
		subject := newHealthcheckStatus()
		result := toDTO(subject)

		assert.NotNil(t, result)
		assert.True(t, result.Healthy)
		assert.Len(t, result.Details, 2)
		assert.Equal(t, "db", result.Details[0].Component)
		assert.True(t, result.Details[0].Healthy)
		assert.Equal(t, "nginx", result.Details[1].Component)
		assert.True(t, result.Details[1].Healthy)
	})

	t.Run("converts unhealthy status to DTO", func(t *testing.T) {
		subject := newHealthcheckStatus()
		subject.Healthy = false
		subject.Details[0].Error = errors.New("db error")

		result := toDTO(subject)

		assert.NotNil(t, result)
		assert.False(t, result.Healthy)
		assert.Len(t, result.Details, 2)
		assert.Equal(t, "db", result.Details[0].Component)
		assert.False(t, result.Details[0].Healthy)
		assert.Equal(t, "nginx", result.Details[1].Component)
		assert.True(t, result.Details[1].Healthy)
	})
}
