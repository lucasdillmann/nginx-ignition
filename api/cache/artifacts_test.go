package cache

import (
	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/cache"
	"dillmann.com.br/nginx-ignition/core/common/pagination"
	"dillmann.com.br/nginx-ignition/core/common/ptr"
)

func newCacheRequestDTO() cacheRequestDTO {
	return cacheRequestDTO{
		Name:            "Static Assets",
		StoragePath:     ptr.Of("/var/cache/nginx"),
		InactiveSeconds: ptr.Of(600),
		MaximumSizeMB:   ptr.Of(1024),
		ConcurrencyLock: concurrencyLockDTO{
			Enabled:        true,
			TimeoutSeconds: ptr.Of(5),
			AgeSeconds:     ptr.Of(5),
		},
		UseStale: []cache.UseStaleOption{cache.ErrorUseStale},
		AllowedMethods: []cache.Method{
			cache.GetMethod,
			cache.HeadMethod,
		},
		BypassRules:              []string{"$cookie_nocache"},
		NoCacheRules:             []string{"$cookie_nocache"},
		FileExtensions:           []string{"jpg", "png"},
		MinimumUsesBeforeCaching: 1,
		BackgroundUpdate:         true,
		Revalidate:               true,
		Durations: []durationDTO{
			{
				StatusCodes:      []string{"200", "302"},
				ValidTimeSeconds: 3600,
			},
		},
	}
}

func newCache() *cache.Cache {
	return &cache.Cache{
		ID:              uuid.New(),
		Name:            "Static Assets",
		StoragePath:     ptr.Of("/var/cache/nginx"),
		InactiveSeconds: ptr.Of(600),
		MaximumSizeMB:   ptr.Of(1024),
		ConcurrencyLock: cache.ConcurrencyLock{
			Enabled:        true,
			TimeoutSeconds: ptr.Of(5),
			AgeSeconds:     ptr.Of(5),
		},
		UseStale: []cache.UseStaleOption{cache.ErrorUseStale},
		AllowedMethods: []cache.Method{
			cache.GetMethod,
			cache.HeadMethod,
		},
		BypassRules:              []string{"$cookie_nocache"},
		NoCacheRules:             []string{"$cookie_nocache"},
		FileExtensions:           []string{"jpg", "png"},
		MinimumUsesBeforeCaching: 1,
		BackgroundUpdate:         true,
		Revalidate:               true,
		Durations: []cache.Duration{
			{
				StatusCodes:      []string{"200", "302"},
				ValidTimeSeconds: 3600,
			},
		},
	}
}

func newCachePage() *pagination.Page[cache.Cache] {
	return pagination.Of([]cache.Cache{
		*newCache(),
	})
}
