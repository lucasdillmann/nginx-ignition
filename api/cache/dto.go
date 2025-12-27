package cache

import (
	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/cache"
)

type cacheRequestDto struct {
	StoragePath              *string                `json:"storagePath"`
	InactiveSeconds          *int                   `json:"inactiveSeconds"`
	MaximumSizeMB            *int                   `json:"maximumSizeMb"`
	ConcurrencyLock          concurrencyLockDto     `json:"concurrencyLock"`
	Name                     string                 `json:"name"`
	UseStale                 []cache.UseStaleOption `json:"useStale"`
	AllowedMethods           []cache.Method         `json:"allowedMethods"`
	BypassRules              []string               `json:"bypassRules"`
	NoCacheRules             []string               `json:"noCacheRules"`
	Durations                []durationDto          `json:"durations"`
	MinimumUsesBeforeCaching int                    `json:"minimumUsesBeforeCaching"`
	BackgroundUpdate         bool                   `json:"backgroundUpdate"`
	Revalidate               bool                   `json:"revalidate"`
}

type cacheResponseDto struct {
	InactiveSeconds          *int                   `json:"inactiveSeconds"`
	StoragePath              *string                `json:"storagePath"`
	MaximumSizeMB            *int                   `json:"maximumSizeMb"`
	ConcurrencyLock          concurrencyLockDto     `json:"concurrencyLock"`
	Name                     string                 `json:"name"`
	UseStale                 []cache.UseStaleOption `json:"useStale"`
	AllowedMethods           []cache.Method         `json:"allowedMethods"`
	BypassRules              []string               `json:"bypassRules"`
	NoCacheRules             []string               `json:"noCacheRules"`
	Durations                []durationDto          `json:"durations"`
	MinimumUsesBeforeCaching int                    `json:"minimumUsesBeforeCaching"`
	ID                       uuid.UUID              `json:"id"`
	Revalidate               bool                   `json:"revalidate"`
	BackgroundUpdate         bool                   `json:"backgroundUpdate"`
}

type concurrencyLockDto struct {
	TimeoutSeconds *int `json:"timeoutSeconds"`
	AgeSeconds     *int `json:"ageSeconds"`
	Enabled        bool `json:"enabled"`
}

type durationDto struct {
	StatusCodes      []int `json:"statusCodes"`
	ValidTimeSeconds int   `json:"validTimeSeconds"`
}
