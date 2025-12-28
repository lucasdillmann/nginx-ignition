package cache

import (
	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/cache"
)

type cacheRequestDTO struct {
	StoragePath                      *string                `json:"storagePath"`
	InactiveSeconds                  *int                   `json:"inactiveSeconds"`
	MaximumSizeMB                    *int                   `json:"maximumSizeMb"`
	ConcurrencyLock                  concurrencyLockDTO     `json:"concurrencyLock"`
	Name                             string                 `json:"name"`
	UseStale                         []cache.UseStaleOption `json:"useStale"`
	AllowedMethods                   []cache.Method         `json:"allowedMethods"`
	BypassRules                      []string               `json:"bypassRules"`
	NoCacheRules                     []string               `json:"noCacheRules"`
	FileExtensions                   []string               `json:"fileExtensions"`
	Durations                        []durationDTO          `json:"durations"`
	MinimumUsesBeforeCaching         int                    `json:"minimumUsesBeforeCaching"`
	BackgroundUpdate                 bool                   `json:"backgroundUpdate"`
	Revalidate                       bool                   `json:"revalidate"`
	IgnoreUpstreamCacheHeaders       bool                   `json:"ignoreUpstreamCacheHeaders"`
	CacheStatusResponseHeaderEnabled bool                   `json:"cacheStatusResponseHeaderEnabled"`
}

type cacheResponseDTO struct {
	InactiveSeconds                  *int                   `json:"inactiveSeconds"`
	StoragePath                      *string                `json:"storagePath"`
	MaximumSizeMB                    *int                   `json:"maximumSizeMb"`
	ConcurrencyLock                  concurrencyLockDTO     `json:"concurrencyLock"`
	Name                             string                 `json:"name"`
	UseStale                         []cache.UseStaleOption `json:"useStale"`
	AllowedMethods                   []cache.Method         `json:"allowedMethods"`
	BypassRules                      []string               `json:"bypassRules"`
	NoCacheRules                     []string               `json:"noCacheRules"`
	FileExtensions                   []string               `json:"fileExtensions"`
	Durations                        []durationDTO          `json:"durations"`
	MinimumUsesBeforeCaching         int                    `json:"minimumUsesBeforeCaching"`
	ID                               uuid.UUID              `json:"id"`
	Revalidate                       bool                   `json:"revalidate"`
	BackgroundUpdate                 bool                   `json:"backgroundUpdate"`
	IgnoreUpstreamCacheHeaders       bool                   `json:"ignoreUpstreamCacheHeaders"`
	CacheStatusResponseHeaderEnabled bool                   `json:"cacheStatusResponseHeaderEnabled"`
}

type concurrencyLockDTO struct {
	TimeoutSeconds *int `json:"timeoutSeconds"`
	AgeSeconds     *int `json:"ageSeconds"`
	Enabled        bool `json:"enabled"`
}

type durationDTO struct {
	StatusCodes      []string `json:"statusCodes"`
	ValidTimeSeconds int      `json:"validTimeSeconds"`
}
