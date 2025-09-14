package client

import (
	"sync"
	"time"

	"github.com/patrickmn/go-cache"
)

type apiCache[T any] struct {
	delegate *cache.Cache
	lock     sync.Mutex
}

func newCache[T any](timeoutSeconds int) *apiCache[T] {
	return &apiCache[T]{
		delegate: cache.New(
			time.Duration(timeoutSeconds)*time.Second,
			time.Duration(timeoutSeconds)*time.Second,
		),
	}
}

func (c *apiCache[T]) get(key string, missProvider func() T) T {
	c.lock.Lock()
	defer c.lock.Unlock()

	if value, found := c.delegate.Get(key); found {
		return value.(T)
	}

	value := missProvider()
	c.delegate.Set(key, value, cache.DefaultExpiration)

	return value
}
