package cache

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"dillmann.com.br/nginx-ignition/core/cache"
	"dillmann.com.br/nginx-ignition/core/common/ptr"
)

func Test_Converter(t *testing.T) {
	t.Run("toDomain", func(t *testing.T) {
		t.Run("converts DTO to domain object", func(t *testing.T) {
			id := uuid.New()
			dto := &cacheRequestDTO{
				Name:            "Test Cache",
				StoragePath:     ptr.Of("/tmp/cache"),
				InactiveSeconds: ptr.Of(600),
				MaximumSizeMB:   ptr.Of(512),
				AllowedMethods:  []cache.Method{cache.GetMethod, cache.HeadMethod},
				UseStale:        []cache.UseStaleOption{cache.ErrorUseStale, cache.TimeoutUseStale},
				BypassRules:     []string{"$cookie_nocache"},
				NoCacheRules:    []string{"$arg_nocache"},
				FileExtensions:  []string{"jpg", "png"},
				Durations: []durationDTO{
					{
						StatusCodes:      []string{"200", "302"},
						ValidTimeSeconds: 3600,
					},
				},
				ConcurrencyLock: concurrencyLockDTO{
					Enabled:        true,
					TimeoutSeconds: ptr.Of(5),
					AgeSeconds:     ptr.Of(10),
				},
			}

			result := toDomain(id, dto)

			assert.NotNil(t, result)
			assert.Equal(t, id, result.ID)
			assert.Equal(t, dto.Name, result.Name)
			assert.Equal(t, dto.StoragePath, result.StoragePath)
			assert.Equal(t, dto.InactiveSeconds, result.InactiveSeconds)
			assert.Equal(t, dto.MaximumSizeMB, result.MaximumSizeMB)
			assert.Equal(t, dto.AllowedMethods, result.AllowedMethods)
			assert.Len(t, result.Durations, 1)
			assert.Equal(t, dto.Durations[0].StatusCodes, result.Durations[0].StatusCodes)
			assert.Equal(t, dto.ConcurrencyLock.Enabled, result.ConcurrencyLock.Enabled)
		})
	})

	t.Run("toResponseDTO", func(t *testing.T) {
		t.Run("converts domain object to response DTO", func(t *testing.T) {
			id := uuid.New()
			domain := &cache.Cache{
				ID:              id,
				Name:            "Test Cache",
				InactiveSeconds: ptr.Of(600),
				Durations: []cache.Duration{
					{
						StatusCodes:      []string{"200"},
						ValidTimeSeconds: 3600,
					},
				},
				ConcurrencyLock: cache.ConcurrencyLock{
					Enabled: true,
				},
			}

			result := toResponseDTO(domain)

			assert.Equal(t, domain.ID, result.ID)
			assert.Equal(t, domain.Name, result.Name)
			assert.Equal(t, domain.InactiveSeconds, result.InactiveSeconds)
			assert.Len(t, result.Durations, 1)
			assert.Equal(t, domain.Durations[0].StatusCodes, result.Durations[0].StatusCodes)
			assert.Equal(t, domain.ConcurrencyLock.Enabled, result.ConcurrencyLock.Enabled)
		})
	})
}
