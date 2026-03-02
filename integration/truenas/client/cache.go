package client

import (
	"sync"
	"time"

	"github.com/patrickmn/go-cache"

	"dillmann.com.br/nginx-ignition/core/common/configuration"
)

var (
	cacheDelegate *cache.Cache
	cacheInitOnce = &sync.Once{}
	cacheInitLock = &sync.Mutex{}
)

func initCache(cfg *configuration.Configuration) error {
	cacheInitLock.Lock()
	defer cacheInitLock.Unlock()

	var err error
	cacheInitOnce.Do(func() {
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
		cacheInitOnce = &sync.Once{}
	}

	return err
}

func getFromCache[T any](key string, missProvider func() (*T, error)) (*T, error) {
	if value, found := cacheDelegate.Get(key); found {
		return value.(*T), nil
	}

	value, err := missProvider()
	if err != nil {
		return nil, err
	}

	cacheDelegate.Set(key, value, cache.DefaultExpiration)
	return value, nil
}
