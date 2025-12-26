package cache

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type cacheModel struct {
	bun.BaseModel `bun:"cache"`

	ID                            uuid.UUID       `bun:"id,pk"`
	Name                          string          `bun:"name,notnull"`
	StoragePath                   *string         `bun:"storage_path"`
	InactiveSeconds               *int            `bun:"inactive_seconds"`
	MaxSizeMB                     *int            `bun:"max_size_mb"`
	AllowedMethods                []string        `bun:"allowed_methods,array,notnull"`
	MinimumUsesBeforeCaching      *int            `bun:"minimum_uses_before_caching"`
	UseStale                      []string        `bun:"use_stale,array,notnull"`
	BackgroundUpdate              *bool           `bun:"background_update"`
	ConcurrencyLockEnabled        bool            `bun:"concurrency_lock_enabled,notnull"`
	ConcurrencyLockTimeoutSeconds *int            `bun:"concurrency_lock_timeout_seconds"`
	ConcurrencyLockAgeSeconds     *int            `bun:"concurrency_lock_age_seconds"`
	Revalidate                    *bool           `bun:"revalidate"`
	BypassRules                   []string        `bun:"bypass_rules,array,notnull"`
	NoCacheRules                  []string        `bun:"no_cache_rules,array,notnull"`
	Durations                     []durationModel `bun:"rel:has-many,join:id=cache_id"`
}

type durationModel struct {
	bun.BaseModel `bun:"cache_duration"`

	ID               uuid.UUID `bun:"id,pk"`
	CacheID          uuid.UUID `bun:"cache_id,notnull"`
	StatusCodes      []int     `bun:"status_codes,array,notnull"`
	ValidTimeSeconds int       `bun:"valid_time_seconds,notnull"`
}
