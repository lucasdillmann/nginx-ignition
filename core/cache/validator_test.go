package cache

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func validCache() *Cache {
	return &Cache{
		Name:                     "test",
		MinimumUsesBeforeCaching: 1,
	}
}

func Test_Validator(t *testing.T) {
	t.Run("validate", func(t *testing.T) {
		t.Run("valid cache passes", func(t *testing.T) {
			cache := validCache()
			val := newValidator()

			err := val.validate(cache)

			assert.NoError(t, err)
		})

		t.Run("empty name fails", func(t *testing.T) {
			cache := validCache()
			cache.Name = ""
			val := newValidator()

			err := val.validate(cache)

			assert.Error(t, err)
		})

		t.Run("whitespace-only name fails", func(t *testing.T) {
			cache := validCache()
			cache.Name = "   "
			val := newValidator()

			err := val.validate(cache)

			assert.Error(t, err)
		})

		t.Run("minimum uses before caching less than 1 fails", func(t *testing.T) {
			cache := validCache()
			cache.MinimumUsesBeforeCaching = 0
			val := newValidator()

			err := val.validate(cache)

			assert.Error(t, err)
		})

		t.Run("inactive seconds less than 1 fails", func(t *testing.T) {
			cache := validCache()
			zero := 0
			cache.InactiveSeconds = &zero
			val := newValidator()

			err := val.validate(cache)

			assert.Error(t, err)
		})

		t.Run("maximum size MB less than 1 fails", func(t *testing.T) {
			cache := validCache()
			zero := 0
			cache.MaximumSizeMB = &zero
			val := newValidator()

			err := val.validate(cache)

			assert.Error(t, err)
		})

		t.Run("relative storage path fails", func(t *testing.T) {
			cache := validCache()
			relativePath := "relative/path"
			cache.StoragePath = &relativePath
			val := newValidator()

			err := val.validate(cache)

			assert.Error(t, err)
		})

		t.Run("absolute storage path passes", func(t *testing.T) {
			cache := validCache()
			absolutePath := "/absolute/path"
			cache.StoragePath = &absolutePath
			val := newValidator()

			err := val.validate(cache)

			assert.NoError(t, err)
		})

		t.Run("concurrency lock enabled without timeout fails", func(t *testing.T) {
			cache := validCache()
			cache.ConcurrencyLock = ConcurrencyLock{
				Enabled: true,
			}
			val := newValidator()

			err := val.validate(cache)

			assert.Error(t, err)
		})

		t.Run("concurrency lock enabled without age fails", func(t *testing.T) {
			cache := validCache()
			timeout := 10
			cache.ConcurrencyLock = ConcurrencyLock{
				Enabled:        true,
				TimeoutSeconds: &timeout,
			}
			val := newValidator()

			err := val.validate(cache)

			assert.Error(t, err)
		})

		t.Run("concurrency lock timeout less than 1 fails", func(t *testing.T) {
			cache := validCache()
			timeout := 0
			age := 5
			cache.ConcurrencyLock = ConcurrencyLock{
				Enabled:        true,
				TimeoutSeconds: &timeout,
				AgeSeconds:     &age,
			}
			val := newValidator()

			err := val.validate(cache)

			assert.Error(t, err)
		})

		t.Run("concurrency lock age less than 1 fails", func(t *testing.T) {
			cache := validCache()
			timeout := 10
			age := 0
			cache.ConcurrencyLock = ConcurrencyLock{
				Enabled:        true,
				TimeoutSeconds: &timeout,
				AgeSeconds:     &age,
			}
			val := newValidator()

			err := val.validate(cache)

			assert.Error(t, err)
		})

		t.Run("concurrency lock enabled with valid values passes", func(t *testing.T) {
			cache := validCache()
			timeout := 10
			age := 5
			cache.ConcurrencyLock = ConcurrencyLock{
				Enabled:        true,
				TimeoutSeconds: &timeout,
				AgeSeconds:     &age,
			}
			val := newValidator()

			err := val.validate(cache)

			assert.NoError(t, err)
		})

		t.Run("invalid HTTP method fails", func(t *testing.T) {
			cache := validCache()
			cache.AllowedMethods = []Method{"INVALID"}
			val := newValidator()

			err := val.validate(cache)

			assert.Error(t, err)
		})

		t.Run("valid HTTP methods pass", func(t *testing.T) {
			cache := validCache()
			cache.AllowedMethods = []Method{GetMethod, PostMethod, PutMethod}
			val := newValidator()

			err := val.validate(cache)

			assert.NoError(t, err)
		})

		t.Run("invalid use stale option fails", func(t *testing.T) {
			cache := validCache()
			cache.UseStale = []UseStaleOption{"INVALID"}
			val := newValidator()

			err := val.validate(cache)

			assert.Error(t, err)
		})

		t.Run("valid use stale options pass", func(t *testing.T) {
			cache := validCache()
			cache.UseStale = []UseStaleOption{ErrorUseStale, TimeoutUseStale}
			val := newValidator()

			err := val.validate(cache)

			assert.NoError(t, err)
		})

		t.Run("duration without status codes fails", func(t *testing.T) {
			cache := validCache()
			cache.Durations = []Duration{
				{
					StatusCodes:      []string{},
					ValidTimeSeconds: 60,
				},
			}
			val := newValidator()

			err := val.validate(cache)

			assert.Error(t, err)
		})

		t.Run("duration with invalid status code fails", func(t *testing.T) {
			cache := validCache()
			cache.Durations = []Duration{
				{
					StatusCodes:      []string{"999"},
					ValidTimeSeconds: 60,
				},
			}
			val := newValidator()

			err := val.validate(cache)

			assert.Error(t, err)
		})

		t.Run("duration status code below range fails", func(t *testing.T) {
			cache := validCache()
			cache.Durations = []Duration{
				{
					StatusCodes:      []string{"99"},
					ValidTimeSeconds: 60,
				},
			}
			val := newValidator()

			err := val.validate(cache)

			assert.Error(t, err)
		})

		t.Run("duration status code above range fails", func(t *testing.T) {
			cache := validCache()
			cache.Durations = []Duration{
				{
					StatusCodes:      []string{"600"},
					ValidTimeSeconds: 60,
				},
			}
			val := newValidator()

			err := val.validate(cache)

			assert.Error(t, err)
		})

		t.Run("duration status code non-integer fails", func(t *testing.T) {
			cache := validCache()
			cache.Durations = []Duration{
				{
					StatusCodes:      []string{"abc"},
					ValidTimeSeconds: 60,
				},
			}
			val := newValidator()

			err := val.validate(cache)

			assert.Error(t, err)
		})

		t.Run("duration valid time seconds less than 1 fails", func(t *testing.T) {
			cache := validCache()
			cache.Durations = []Duration{
				{
					StatusCodes:      []string{"200"},
					ValidTimeSeconds: 0,
				},
			}
			val := newValidator()

			err := val.validate(cache)

			assert.Error(t, err)
		})

		t.Run("duration with valid status code passes", func(t *testing.T) {
			cache := validCache()
			cache.Durations = []Duration{
				{
					StatusCodes:      []string{"200", "404"},
					ValidTimeSeconds: 60,
				},
			}
			val := newValidator()

			err := val.validate(cache)

			assert.NoError(t, err)
		})

		t.Run("empty file extension fails", func(t *testing.T) {
			cache := validCache()
			cache.FileExtensions = []string{""}
			val := newValidator()

			err := val.validate(cache)

			assert.Error(t, err)
		})

		t.Run("whitespace-only file extension fails", func(t *testing.T) {
			cache := validCache()
			cache.FileExtensions = []string{"   "}
			val := newValidator()

			err := val.validate(cache)

			assert.Error(t, err)
		})

		t.Run("file extension starting with dot fails", func(t *testing.T) {
			cache := validCache()
			cache.FileExtensions = []string{".txt"}
			val := newValidator()

			err := val.validate(cache)

			assert.Error(t, err)
		})

		t.Run("file extension without dot passes", func(t *testing.T) {
			cache := validCache()
			cache.FileExtensions = []string{"txt", "jpg"}
			val := newValidator()

			err := val.validate(cache)

			assert.NoError(t, err)
		})
	})
}
