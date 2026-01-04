package settings

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_toDTO(t *testing.T) {
	t.Run("returns nil when input is nil", func(t *testing.T) {
		result := toDTO(nil)
		assert.Nil(t, result)
	})

	t.Run("converts domain object to DTO", func(t *testing.T) {
		subject := newSettings()
		result := toDTO(subject)

		assert.NotNil(t, result)
		assert.Equal(t, subject.Nginx.GzipEnabled, *result.Nginx.GzipEnabled)
	})
}

func Test_toDomain(t *testing.T) {
	t.Run("returns nil when critical fields are missing", func(t *testing.T) {
		result := toDomain(&settingsDTO{})
		assert.Nil(t, result)
	})

	t.Run("converts DTO to domain object", func(t *testing.T) {
		payload := newSettingsDTO()
		result := toDomain(payload)

		assert.NotNil(t, result)
		assert.Equal(t, *payload.Nginx.GzipEnabled, result.Nginx.GzipEnabled)
	})
}
