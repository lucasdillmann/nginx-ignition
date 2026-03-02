package client

import (
	"sync"
	"time"

	"github.com/patrickmn/go-cache"

	"dillmann.com.br/nginx-ignition/core/common/configuration"
)

var (
	cacheDelegate *cache.Cache
	cacheInitLock = &sync.Once{}
	cacheMainLock = &sync.Mutex{}
)

func initCache(cfg *configuration.Configuration) error {
	var err error

	cacheInitLock.Do(func() {
		cacheDuration, innerErr := cfg.GetInt(
			"nginx-ignition.integration.truenas.api-cache-timeout-seconds",
		)
		if innerErr != nil {
			err = innerErr
			return
		}

		cacheDelegate = cache.New(
			time.Duration(cacheDuration)*time.Second,
			time.Duration(cacheDuration)*time.Second,
		)
	})

	if err != nil {
		cacheInitLock = &sync.Once{}
	}

	return err
}

func getFromCache[T any](key string, missProvider func() T) T {
	cacheMainLock.Lock()
	defer cacheMainLock.Unlock()

	if value, found := cacheDelegate.Get(key); found {
		return value.(T)
	}

	value := missProvider()
	cacheDelegate.Set(key, value, cache.DefaultExpiration)

	return value
}
