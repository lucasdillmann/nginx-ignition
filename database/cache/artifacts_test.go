package cache

import (
	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/cache"
)

func newCache() *cache.Cache {
	return &cache.Cache{
		ID:          uuid.New(),
		Name:        "Test Cache",
		StoragePath: new("/tmp/tests"),
		ConcurrencyLock: cache.ConcurrencyLock{
			Enabled:        true,
			TimeoutSeconds: new(5),
			AgeSeconds:     new(10),
		},
		AllowedMethods: []cache.Method{
			cache.GetMethod,
			cache.HeadMethod,
		},
		UseStale: []cache.UseStaleOption{
			cache.HTTP500UseStale,
			cache.UpdatingUseStale,
		},
		Durations: []cache.Duration{
			{
				StatusCodes:      []string{"200", "301"},
				ValidTimeSeconds: 60,
			},
			{
				StatusCodes:      []string{"404"},
				ValidTimeSeconds: 10,
			},
		},
		FileExtensions: []string{
			"jpg",
			"png",
		},
		BypassRules:              []string{},
		NoCacheRules:             []string{},
		MinimumUsesBeforeCaching: 1,
		InactiveSeconds:          new(300),
		MaximumSizeMB:            new(1024),
	}
}
