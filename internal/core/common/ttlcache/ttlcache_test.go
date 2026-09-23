package ttlcache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Run("returns error for zero ttl", func(t *testing.T) {
		cache, err := New[string, string](0)
		assert.Error(t, err)
		assert.Nil(t, cache)
		assert.Contains(t, err.Error(), "ttl cannot be zero or negative")
	})

	t.Run("returns error for negative ttl", func(t *testing.T) {
		cache, err := New[string, string](-time.Hour)
		assert.Error(t, err)
		assert.Nil(t, cache)
		assert.Contains(t, err.Error(), "ttl cannot be zero or negative")
	})

	t.Run("returns cache for positive ttl", func(t *testing.T) {
		cache, err := New[string, string](time.Hour)
		require.NoError(t, err)
		assert.NotNil(t, cache)
	})
}

func TestCache_SetAndGet(t *testing.T) {
	cache, err := New[string, string](time.Hour)
	require.NoError(t, err)

	cache.Set("key1", "value1")
	val, ok := cache.Get("key1")

	require.True(t, ok)
	assert.Equal(t, "value1", val)
}

func TestCache_GetMissingKey(t *testing.T) {
	cache, err := New[string, string](time.Hour)
	require.NoError(t, err)

	val, ok := cache.Get("missing")
	assert.False(t, ok)
	assert.Empty(t, val)
}

func TestCache_Delete(t *testing.T) {
	cache, err := New[string, string](time.Hour)
	require.NoError(t, err)

	cache.Set("key1", "value1")
	cache.Delete("key1")
	val, ok := cache.Get("key1")

	assert.False(t, ok)
	assert.Empty(t, val)
}

func TestCache_Expiration(t *testing.T) {
	cache, err := New[string, string](10 * time.Millisecond)
	require.NoError(t, err)

	cache.Set("key1", "value1")
	val, ok := cache.Get("key1")
	require.True(t, ok)
	assert.Equal(t, "value1", val)

	time.Sleep(20 * time.Millisecond)

	val, ok = cache.Get("key1")
	assert.False(t, ok)
	assert.Empty(t, val)
}

func TestCache_GenericTypes(t *testing.T) {
	cache, err := New[int, []string](time.Hour)
	require.NoError(t, err)

	cache.Set(42, []string{"a", "b"})
	val, ok := cache.Get(42)
	require.True(t, ok)
	assert.Equal(t, []string{"a", "b"}, val)
}

func TestCache_CleanupLoop(t *testing.T) {
	cache, err := New[string, string](10 * time.Millisecond)
	require.NoError(t, err)

	cache.Set("key1", "value1")
	cache.Set("key2", "value2")

	time.Sleep(25 * time.Millisecond)

	_, ok := cache.Get("key1")
	assert.False(t, ok)

	_, ok = cache.Get("key2")
	assert.False(t, ok)
}

func TestCache_ConcurrentAccess(t *testing.T) {
	cache, err := New[int, int](time.Hour)
	require.NoError(t, err)

	done := make(chan bool, 100)
	for index := range 50 {
		go func(n int) {
			cache.Set(n, n*2)
			done <- true
		}(index)

		go func(n int) {
			_, _ = cache.Get(n)
			done <- true
		}(index)
	}

	for range 100 {
		<-done
	}

	require.Equal(t, 50, len(cache.items))
	for index := range 50 {
		val, ok := cache.Get(index)

		require.True(t, ok)
		assert.Equal(t, index*2, val)
	}
}
