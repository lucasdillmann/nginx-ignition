package cache

import (
	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/cache"
	"dillmann.com.br/nginx-ignition/core/common/ptr"
)

func newCache() *cache.Cache {
	return &cache.Cache{
		ID:          uuid.New(),
		Name:        "Test Cache",
		StoragePath: ptr.Of("/tmp/tests"),
		ConcurrencyLock: cache.ConcurrencyLock{
			Enabled:        true,
			TimeoutSeconds: ptr.Of(5),
			AgeSeconds:     ptr.Of(10),
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
		InactiveSeconds:          ptr.Of(300),
		MaximumSizeMB:            ptr.Of(1024),
	}
}
