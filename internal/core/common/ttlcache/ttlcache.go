package ttlcache

import (
	"sync"
	"time"
)

type Cache[K comparable, V any] struct {
	items map[K]entry[V]
	ttl   time.Duration
	mu    sync.RWMutex
}

type entry[V any] struct {
	value      V
	expiration time.Time
}

func New[K comparable, V any](ttl time.Duration) *Cache[K, V] {
	cache := &Cache[K, V]{
		items: make(map[K]entry[V]),
		ttl:   ttl,
	}

	go cache.cleanupLoop()

	return cache
}

func (c *Cache[K, V]) Set(key K, value V) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = entry[V]{
		value:      value,
		expiration: time.Now().Add(c.ttl),
	}
}

func (c *Cache[K, V]) Get(key K) (V, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	e, ok := c.items[key]
	if !ok {
		var zero V
		return zero, false
	}

	if time.Now().After(e.expiration) {
		var zero V
		return zero, false
	}

	return e.value, true
}

func (c *Cache[K, V]) Delete(key K) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.items, key)
}

func (c *Cache[K, V]) cleanupLoop() {
	ticker := time.NewTicker(c.ttl / 2)
	defer ticker.Stop()

	for range ticker.C {
		c.mu.Lock()

		now := time.Now()
		for k, v := range c.items {
			if now.After(v.expiration) {
				delete(c.items, k)
			}
		}

		c.mu.Unlock()
	}
}
