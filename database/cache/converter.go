package cache

import (
	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/cache"
)

func toDomain(model *cacheModel) *cache.Cache {
	durations := make([]cache.Duration, len(model.Durations))
	for index, duration := range model.Durations {
		durations[index] = cache.Duration{
			StatusCodes:      duration.StatusCodes,
			ValidTimeSeconds: duration.ValidTimeSeconds,
		}
	}

	allowedMethods := make([]cache.Method, len(model.AllowedMethods))
	for index, method := range model.AllowedMethods {
		allowedMethods[index] = cache.Method(method)
	}

	useStale := make([]cache.UseStaleOption, len(model.UseStale))
	for index, option := range model.UseStale {
		useStale[index] = cache.UseStaleOption(option)
	}

	return &cache.Cache{
		ID:                       model.ID,
		Name:                     model.Name,
		StoragePath:              model.StoragePath,
		InactiveSeconds:          model.InactiveSeconds,
		MaxSizeMB:                model.MaxSizeMB,
		AllowedMethods:           allowedMethods,
		MinimumUsesBeforeCaching: model.MinimumUsesBeforeCaching,
		UseStale:                 useStale,
		BackgroundUpdate:         model.BackgroundUpdate,
		ConcurrencyLock: cache.ConcurrencyLock{
			Enabled:        model.ConcurrencyLockEnabled,
			TimeoutSeconds: model.ConcurrencyLockTimeoutSeconds,
			AgeSeconds:     model.ConcurrencyLockAgeSeconds,
		},
		Revalidate:   model.Revalidate,
		BypassRules:  model.BypassRules,
		NoCacheRules: model.NoCacheRules,
		Durations:    durations,
	}
}

func toModel(domain *cache.Cache) *cacheModel {
	durations := make([]*durationModel, len(domain.Durations))
	for index, duration := range domain.Durations {
		durations[index] = &durationModel{
			ID:               uuid.New(),
			CacheID:          domain.ID,
			StatusCodes:      duration.StatusCodes,
			ValidTimeSeconds: duration.ValidTimeSeconds,
		}
	}

	allowedMethods := make([]string, len(domain.AllowedMethods))
	for index, method := range domain.AllowedMethods {
		allowedMethods[index] = string(method)
	}

	useStale := make([]string, len(domain.UseStale))
	for index, option := range domain.UseStale {
		useStale[index] = string(option)
	}

	return &cacheModel{
		ID:                            domain.ID,
		Name:                          domain.Name,
		StoragePath:                   domain.StoragePath,
		InactiveSeconds:               domain.InactiveSeconds,
		MaxSizeMB:                     domain.MaxSizeMB,
		AllowedMethods:                allowedMethods,
		MinimumUsesBeforeCaching:      domain.MinimumUsesBeforeCaching,
		UseStale:                      useStale,
		BackgroundUpdate:              domain.BackgroundUpdate,
		ConcurrencyLockEnabled:        domain.ConcurrencyLock.Enabled,
		ConcurrencyLockTimeoutSeconds: domain.ConcurrencyLock.TimeoutSeconds,
		ConcurrencyLockAgeSeconds:     domain.ConcurrencyLock.AgeSeconds,
		Revalidate:                    domain.Revalidate,
		BypassRules:                   domain.BypassRules,
		NoCacheRules:                  domain.NoCacheRules,
		Durations:                     durations,
	}
}
