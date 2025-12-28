package cache

import (
	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/cache"
)

func toDomain(id uuid.UUID, dto *cacheRequestDTO) *cache.Cache {
	durations := make([]cache.Duration, len(dto.Durations))
	for index, duration := range dto.Durations {
		durations[index] = cache.Duration{
			StatusCodes:      duration.StatusCodes,
			ValidTimeSeconds: duration.ValidTimeSeconds,
		}
	}

	return &cache.Cache{
		ID:                               id,
		Name:                             dto.Name,
		StoragePath:                      dto.StoragePath,
		InactiveSeconds:                  dto.InactiveSeconds,
		MaximumSizeMB:                    dto.MaximumSizeMB,
		AllowedMethods:                   dto.AllowedMethods,
		MinimumUsesBeforeCaching:         dto.MinimumUsesBeforeCaching,
		UseStale:                         dto.UseStale,
		BackgroundUpdate:                 dto.BackgroundUpdate,
		Revalidate:                       dto.Revalidate,
		BypassRules:                      dto.BypassRules,
		NoCacheRules:                     dto.NoCacheRules,
		FileExtensions:                   dto.FileExtensions,
		IgnoreUpstreamCacheHeaders:       dto.IgnoreUpstreamCacheHeaders,
		CacheStatusResponseHeaderEnabled: dto.CacheStatusResponseHeaderEnabled,
		Durations:                        durations,
		ConcurrencyLock: cache.ConcurrencyLock{
			Enabled:        dto.ConcurrencyLock.Enabled,
			TimeoutSeconds: dto.ConcurrencyLock.TimeoutSeconds,
			AgeSeconds:     dto.ConcurrencyLock.AgeSeconds,
		},
	}
}

func toResponseDTO(domain *cache.Cache) cacheResponseDTO {
	durations := make([]durationDTO, len(domain.Durations))
	for index, duration := range domain.Durations {
		durations[index] = durationDTO{
			StatusCodes:      duration.StatusCodes,
			ValidTimeSeconds: duration.ValidTimeSeconds,
		}
	}

	return cacheResponseDTO{
		ID:                               domain.ID,
		Name:                             domain.Name,
		StoragePath:                      domain.StoragePath,
		InactiveSeconds:                  domain.InactiveSeconds,
		MaximumSizeMB:                    domain.MaximumSizeMB,
		AllowedMethods:                   domain.AllowedMethods,
		MinimumUsesBeforeCaching:         domain.MinimumUsesBeforeCaching,
		UseStale:                         domain.UseStale,
		BackgroundUpdate:                 domain.BackgroundUpdate,
		Revalidate:                       domain.Revalidate,
		BypassRules:                      domain.BypassRules,
		NoCacheRules:                     domain.NoCacheRules,
		FileExtensions:                   domain.FileExtensions,
		IgnoreUpstreamCacheHeaders:       domain.IgnoreUpstreamCacheHeaders,
		CacheStatusResponseHeaderEnabled: domain.CacheStatusResponseHeaderEnabled,
		Durations:                        durations,
		ConcurrencyLock: concurrencyLockDTO{
			Enabled:        domain.ConcurrencyLock.Enabled,
			TimeoutSeconds: domain.ConcurrencyLock.TimeoutSeconds,
			AgeSeconds:     domain.ConcurrencyLock.AgeSeconds,
		},
	}
}
