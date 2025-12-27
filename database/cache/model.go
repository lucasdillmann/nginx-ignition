package cache

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type cacheModel struct {
	bun.BaseModel                 `bun:"cache"`
	ConcurrencyLockAgeSeconds     *int            `bun:"concurrency_lock_age_seconds"`
	ConcurrencyLockTimeoutSeconds *int            `bun:"concurrency_lock_timeout_seconds"`
	StoragePath                   *string         `bun:"storage_path"`
	InactiveSeconds               *int            `bun:"inactive_seconds"`
	MaximumSizeMB                 *int            `bun:"maximum_size_mb"`
	Name                          string          `bun:"name,notnull"`
	BypassRules                   []string        `bun:"bypass_rules,array,notnull"`
	UseStale                      []string        `bun:"use_stale,array,notnull"`
	AllowedMethods                []string        `bun:"allowed_methods,array,notnull"`
	NoCacheRules                  []string        `bun:"no_cache_rules,array,notnull"`
	Durations                     []durationModel `bun:"rel:has-many,join:id=cache_id"`
	MinimumUsesBeforeCaching      int             `bun:"minimum_uses_before_caching,notnull"`
	ID                            uuid.UUID       `bun:"id,pk"`
	BackgroundUpdate              bool            `bun:"background_update,notnull"`
	Revalidate                    bool            `bun:"revalidate,notnull"`
	ConcurrencyLockEnabled        bool            `bun:"concurrency_lock_enabled,notnull"`
}

type durationModel struct {
	bun.BaseModel    `bun:"cache_duration"`
	StatusCodes      []int     `bun:"status_codes,array,notnull"`
	ValidTimeSeconds int       `bun:"valid_time_seconds,notnull"`
	ID               uuid.UUID `bun:"id,pk"`
	CacheID          uuid.UUID `bun:"cache_id,notnull"`
}
