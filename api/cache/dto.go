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
	ID                       uuid.UUID              `json:"id"`
	Name                     string                 `json:"name"`
	StoragePath              *string                `json:"storagePath"`
	InactiveSeconds          *int                   `json:"inactiveSeconds"`
	MaxSizeMB                *int                   `json:"maxSizeMb"`
	AllowedMethods           []cache.Method         `json:"allowedMethods"`
	MinimumUsesBeforeCaching *int                   `json:"minimumUsesBeforeCaching"`
	UseStale                 []cache.UseStaleOption `json:"useStale"`
	BackgroundUpdate         *bool                  `json:"backgroundUpdate"`
	ConcurrencyLock          concurrencyLockDto     `json:"concurrencyLock"`
	Revalidate               *bool                  `json:"revalidate"`
	BypassRules              []string               `json:"bypassRules"`
	NoCacheRules             []string               `json:"noCacheRules"`
	Durations                []durationDto          `json:"durations"`
}

type concurrencyLockDto struct {
	Enabled        bool `json:"enabled"`
	TimeoutSeconds *int `json:"timeoutSeconds"`
	AgeSeconds     *int `json:"ageSeconds"`
}

type durationDto struct {
	StatusCodes      []int `json:"statusCodes"`
	ValidTimeSeconds int   `json:"validTimeSeconds"`
}
