package cache

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_toDomain(t *testing.T) {
	t.Run("converts DTO to domain object", func(t *testing.T) {
		id := uuid.New()
		payload := newCacheRequestDTO()
		result := toDomain(id, &payload)

		assert.NotNil(t, result)
		assert.Equal(t, id, result.ID)
		assert.Equal(t, payload.Name, result.Name)
		assert.Equal(t, payload.StoragePath, result.StoragePath)
		assert.Equal(t, payload.InactiveSeconds, result.InactiveSeconds)
		assert.Equal(t, payload.MaximumSizeMB, result.MaximumSizeMB)
		assert.Equal(t, payload.AllowedMethods, result.AllowedMethods)
		assert.Len(t, result.Durations, 1)
		assert.Equal(t, payload.Durations[0].StatusCodes, result.Durations[0].StatusCodes)
		assert.Equal(t, payload.ConcurrencyLock.Enabled, result.ConcurrencyLock.Enabled)
	})
}

func Test_toResponseDTO(t *testing.T) {
	t.Run("converts domain object to response DTO", func(t *testing.T) {
		subject := newCache()
		result := toResponseDTO(subject)

		assert.Equal(t, subject.ID, result.ID)
		assert.Equal(t, subject.Name, result.Name)
		assert.Equal(t, subject.InactiveSeconds, result.InactiveSeconds)
		assert.Len(t, result.Durations, 1)
		assert.Equal(t, subject.Durations[0].StatusCodes, result.Durations[0].StatusCodes)
		assert.Equal(t, subject.ConcurrencyLock.Enabled, result.ConcurrencyLock.Enabled)
	})
}
