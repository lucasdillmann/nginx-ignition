package cache

import (
	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/cache"
)

type cacheRequestDto struct {
	Name                     *string                `json:"name"`
	StoragePath              *string                `json:"storagePath"`
	InactiveSeconds          *int                   `json:"inactiveSeconds"`
	MaxSizeMB                *int                   `json:"maxSizeMb"`
	AllowedMethods           []cache.Method         `json:"allowedMethods"`
	MinimumUsesBeforeCaching *int                   `json:"minimumUsesBeforeCaching"`
	UseStale                 []cache.UseStaleOption `json:"useStale"`
	BackgroundUpdate         *bool                  `json:"backgroundUpdate"`
	ConcurrencyLock          *concurrencyLockDto    `json:"concurrencyLock"`
	Revalidate               *bool                  `json:"revalidate"`
	BypassRules              []string               `json:"bypassRules"`
	NoCacheRules             []string               `json:"noCacheRules"`
	Durations                []durationDto          `json:"durations"`
}

type cacheResponseDto struct {
	ConcurrencyLock          concurrencyLockDto     `json:"concurrencyLock"`
	InactiveSeconds          *int                   `json:"inactiveSeconds"`
	StoragePath              *string                `json:"storagePath"`
	MaxSizeMB                *int                   `json:"maxSizeMb"`
	MinimumUsesBeforeCaching *int                   `json:"minimumUsesBeforeCaching"`
	BackgroundUpdate         *bool                  `json:"backgroundUpdate"`
	Revalidate               *bool                  `json:"revalidate"`
	Name                     string                 `json:"name"`
	AllowedMethods           []cache.Method         `json:"allowedMethods"`
	UseStale                 []cache.UseStaleOption `json:"useStale"`
	BypassRules              []string               `json:"bypassRules"`
	NoCacheRules             []string               `json:"noCacheRules"`
	Durations                []durationDto          `json:"durations"`
	ID                       uuid.UUID              `json:"id"`
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
