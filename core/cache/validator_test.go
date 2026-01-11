package cache

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_validator(t *testing.T) {
	t.Run("validate", func(t *testing.T) {
		t.Run("valid cache passes", func(t *testing.T) {
			cache := newCache()
			cacheValidator := newValidator()

			err := cacheValidator.validate(t.Context(), cache)

			assert.NoError(t, err)
		})

		t.Run("empty name fails", func(t *testing.T) {
			cache := newCache()
			cache.Name = ""
			cacheValidator := newValidator()

			err := cacheValidator.validate(t.Context(), cache)

			assert.Error(t, err)
		})

		t.Run("whitespace-only name fails", func(t *testing.T) {
			cache := newCache()
			cache.Name = "   "
			cacheValidator := newValidator()

			err := cacheValidator.validate(t.Context(), cache)

			assert.Error(t, err)
		})

		t.Run("minimum uses before caching less than 1 fails", func(t *testing.T) {
			cache := newCache()
			cache.MinimumUsesBeforeCaching = 0
			cacheValidator := newValidator()

			err := cacheValidator.validate(t.Context(), cache)

			assert.Error(t, err)
		})

		t.Run("inactive seconds less than 1 fails", func(t *testing.T) {
			cache := newCache()
			zero := 0
			cache.InactiveSeconds = &zero
			cacheValidator := newValidator()

			err := cacheValidator.validate(t.Context(), cache)

			assert.Error(t, err)
		})

		t.Run("maximum size MB less than 1 fails", func(t *testing.T) {
			cache := newCache()
			zero := 0
			cache.MaximumSizeMB = &zero
			cacheValidator := newValidator()

			err := cacheValidator.validate(t.Context(), cache)

			assert.Error(t, err)
		})

		t.Run("relative storage path fails", func(t *testing.T) {
			cache := newCache()
			relativePath := "relative/path"
			cache.StoragePath = &relativePath
			cacheValidator := newValidator()

			err := cacheValidator.validate(t.Context(), cache)

			assert.Error(t, err)
		})

		t.Run("absolute storage path passes", func(t *testing.T) {
			cache := newCache()
			absolutePath := "/absolute/path"
			cache.StoragePath = &absolutePath
			cacheValidator := newValidator()

			err := cacheValidator.validate(t.Context(), cache)

			assert.NoError(t, err)
		})

		t.Run("concurrency lock enabled without timeout fails", func(t *testing.T) {
			cache := newCache()
			cache.ConcurrencyLock = ConcurrencyLock{
				Enabled: true,
			}
			cacheValidator := newValidator()

			err := cacheValidator.validate(t.Context(), cache)

			assert.Error(t, err)
		})

		t.Run("concurrency lock enabled without age fails", func(t *testing.T) {
			cache := newCache()
			timeout := 10
			cache.ConcurrencyLock = ConcurrencyLock{
				Enabled:        true,
				TimeoutSeconds: &timeout,
			}
			cacheValidator := newValidator()

			err := cacheValidator.validate(t.Context(), cache)

			assert.Error(t, err)
		})

		t.Run("concurrency lock timeout less than 1 fails", func(t *testing.T) {
			cache := newCache()
			timeout := 0
			age := 5
			cache.ConcurrencyLock = ConcurrencyLock{
				Enabled:        true,
				TimeoutSeconds: &timeout,
				AgeSeconds:     &age,
			}
			cacheValidator := newValidator()

			err := cacheValidator.validate(t.Context(), cache)

			assert.Error(t, err)
		})

		t.Run("concurrency lock age less than 1 fails", func(t *testing.T) {
			cache := newCache()
			timeout := 10
			age := 0
			cache.ConcurrencyLock = ConcurrencyLock{
				Enabled:        true,
				TimeoutSeconds: &timeout,
				AgeSeconds:     &age,
			}
			cacheValidator := newValidator()

			err := cacheValidator.validate(t.Context(), cache)

			assert.Error(t, err)
		})

		t.Run("concurrency lock enabled with valid values passes", func(t *testing.T) {
			cache := newCache()
			timeout := 10
			age := 5
			cache.ConcurrencyLock = ConcurrencyLock{
				Enabled:        true,
				TimeoutSeconds: &timeout,
				AgeSeconds:     &age,
			}
			cacheValidator := newValidator()

			err := cacheValidator.validate(t.Context(), cache)

			assert.NoError(t, err)
		})

		t.Run("invalid HTTP method fails", func(t *testing.T) {
			cache := newCache()
			cache.AllowedMethods = []Method{"INVALID"}
			cacheValidator := newValidator()

			err := cacheValidator.validate(t.Context(), cache)

			assert.Error(t, err)
		})

		t.Run("valid HTTP methods pass", func(t *testing.T) {
			cache := newCache()
			cache.AllowedMethods = []Method{GetMethod, PostMethod, PutMethod}
			cacheValidator := newValidator()

			err := cacheValidator.validate(t.Context(), cache)

			assert.NoError(t, err)
		})

		t.Run("invalid use stale option fails", func(t *testing.T) {
			cache := newCache()
			cache.UseStale = []UseStaleOption{"INVALID"}
			cacheValidator := newValidator()

			err := cacheValidator.validate(t.Context(), cache)

			assert.Error(t, err)
		})

		t.Run("valid use stale options pass", func(t *testing.T) {
			cache := newCache()
			cache.UseStale = []UseStaleOption{ErrorUseStale, TimeoutUseStale}
			cacheValidator := newValidator()

			err := cacheValidator.validate(t.Context(), cache)

			assert.NoError(t, err)
		})

		t.Run("duration without status codes fails", func(t *testing.T) {
			cache := newCache()
			cache.Durations = []Duration{
				{
					StatusCodes:      []string{},
					ValidTimeSeconds: 60,
				},
			}
			cacheValidator := newValidator()

			err := cacheValidator.validate(t.Context(), cache)

			assert.Error(t, err)
		})

		t.Run("duration with invalid status code fails", func(t *testing.T) {
			cache := newCache()
			cache.Durations = []Duration{
				{
					StatusCodes:      []string{"999"},
					ValidTimeSeconds: 60,
				},
			}
			cacheValidator := newValidator()

			err := cacheValidator.validate(t.Context(), cache)

			assert.Error(t, err)
		})

		t.Run("duration status code below range fails", func(t *testing.T) {
			cache := newCache()
			cache.Durations = []Duration{
				{
					StatusCodes:      []string{"99"},
					ValidTimeSeconds: 60,
				},
			}
			cacheValidator := newValidator()

			err := cacheValidator.validate(t.Context(), cache)

			assert.Error(t, err)
		})

		t.Run("duration status code above range fails", func(t *testing.T) {
			cache := newCache()
			cache.Durations = []Duration{
				{
					StatusCodes:      []string{"600"},
					ValidTimeSeconds: 60,
				},
			}
			cacheValidator := newValidator()

			err := cacheValidator.validate(t.Context(), cache)

			assert.Error(t, err)
		})

		t.Run("duration status code non-integer fails", func(t *testing.T) {
			cache := newCache()
			cache.Durations = []Duration{
				{
					StatusCodes:      []string{"abc"},
					ValidTimeSeconds: 60,
				},
			}
			cacheValidator := newValidator()

			err := cacheValidator.validate(t.Context(), cache)

			assert.Error(t, err)
		})

		t.Run("duration valid time seconds less than 1 fails", func(t *testing.T) {
			cache := newCache()
			cache.Durations = []Duration{
				{
					StatusCodes:      []string{"200"},
					ValidTimeSeconds: 0,
				},
			}
			cacheValidator := newValidator()

			err := cacheValidator.validate(t.Context(), cache)

			assert.Error(t, err)
		})

		t.Run("duration with valid status code passes", func(t *testing.T) {
			cache := newCache()
			cache.Durations = []Duration{
				{
					StatusCodes:      []string{"200", "404"},
					ValidTimeSeconds: 60,
				},
			}
			cacheValidator := newValidator()

			err := cacheValidator.validate(t.Context(), cache)

			assert.NoError(t, err)
		})

		t.Run("empty file extension fails", func(t *testing.T) {
			cache := newCache()
			cache.FileExtensions = []string{""}
			cacheValidator := newValidator()

			err := cacheValidator.validate(t.Context(), cache)

			assert.Error(t, err)
		})

		t.Run("whitespace-only file extension fails", func(t *testing.T) {
			cache := newCache()
			cache.FileExtensions = []string{"   "}
			cacheValidator := newValidator()

			err := cacheValidator.validate(t.Context(), cache)

			assert.Error(t, err)
		})

		t.Run("file extension starting with dot fails", func(t *testing.T) {
			cache := newCache()
			cache.FileExtensions = []string{".txt"}
			cacheValidator := newValidator()

			err := cacheValidator.validate(t.Context(), cache)

			assert.Error(t, err)
		})

		t.Run("file extension without dot passes", func(t *testing.T) {
			cache := newCache()
			cache.FileExtensions = []string{"txt", "jpg"}
			cacheValidator := newValidator()

			err := cacheValidator.validate(t.Context(), cache)

			assert.NoError(t, err)
		})
	})
}
