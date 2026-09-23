package ttlcache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCache_SetAndGet(t *testing.T) {
	c := New[string, string](time.Hour)
	defer func() { _ = c }()

	c.Set("key1", "value1")
	val, ok := c.Get("key1")
	require.True(t, ok)
	assert.Equal(t, "value1", val)
}

func TestCache_GetMissingKey(t *testing.T) {
	c := New[string, string](time.Hour)
	defer func() { _ = c }()

	val, ok := c.Get("missing")
	assert.False(t, ok)
	assert.Empty(t, val)
}

func TestCache_Delete(t *testing.T) {
	c := New[string, string](time.Hour)
	defer func() { _ = c }()

	c.Set("key1", "value1")
	c.Delete("key1")
	val, ok := c.Get("key1")
	assert.False(t, ok)
	assert.Empty(t, val)
}

func TestCache_Expiration(t *testing.T) {
	c := New[string, string](10 * time.Millisecond)
	defer func() { _ = c }()

	c.Set("key1", "value1")
	val, ok := c.Get("key1")
	require.True(t, ok)
	assert.Equal(t, "value1", val)

	time.Sleep(20 * time.Millisecond)

	val, ok = c.Get("key1")
	assert.False(t, ok)
	assert.Empty(t, val)
}

func TestCache_GenericTypes(t *testing.T) {
	c := New[int, []string](time.Hour)
	defer func() { _ = c }()

	c.Set(42, []string{"a", "b"})
	val, ok := c.Get(42)
	require.True(t, ok)
	assert.Equal(t, []string{"a", "b"}, val)
}

func TestCache_CleanupLoop(t *testing.T) {
	c := New[string, string](10 * time.Millisecond)
	defer func() { _ = c }()

	c.Set("key1", "value1")
	c.Set("key2", "value2")

	time.Sleep(25 * time.Millisecond)

	_, ok := c.Get("key1")
	assert.False(t, ok)

	_, ok = c.Get("key2")
	assert.False(t, ok)
}

func TestCache_ConcurrentAccess(t *testing.T) {
	c := New[int, int](time.Hour)
	defer func() { _ = c }()

	done := make(chan bool, 100)
	for index := range 50 {
		go func(n int) {
			c.Set(n, n*2)
			done <- true
		}(index)

		go func(n int) {
			_, _ = c.Get(n)
			done <- true
		}(index)
	}

	for range 100 {
		<-done
	}

	require.Equal(t, 50, len(c.items))
	for index := range 50 {
		val, ok := c.Get(index)

		require.True(t, ok)
		assert.Equal(t, index*2, val)
	}
}
