package client

import (
	"sync"
	"time"

	"github.com/patrickmn/go-cache"

	"dillmann.com.br/nginx-ignition/core/common/configuration"
)

var (
	cacheDelegate    *cache.Cache
	cacheInitLock    = &sync.Mutex{}
	cacheInitialized = false
)

func initCache(cfg *configuration.Configuration) error {
	cacheInitLock.Lock()
	defer cacheInitLock.Unlock()

	if cacheInitialized {
		return nil
	}

	cacheDuration, err := cfg.GetInt(
		"nginx-ignition.integration.truenas.api-cache-timeout-seconds",
	)
	if err != nil {
		return err
	}

	cacheDelegate = cache.New(
		time.Duration(cacheDuration)*time.Second,
		time.Duration(cacheDuration)*time.Second,
	)

	cacheInitialized = true
	return nil
}

func getFromCache[T any](key string, missProvider func() (*T, error)) (*T, error) {
	if value, found := cacheDelegate.Get(key); found {
		if value == nil {
			return nil, nil
		}

		return value.(*T), nil
	}

	value, err := missProvider()
	if err != nil {
		return nil, err
	}

	cacheDelegate.Set(key, value, cache.DefaultExpiration)
	return value, nil
}
