package cache

import (
	"github.com/google/uuid"
)

func newCache() *Cache {
	return &Cache{
		ID:                       uuid.New(),
		Name:                     "Default Cache",
		MinimumUsesBeforeCaching: 1,
		InactiveSeconds:          nil,
		MaximumSizeMB:            nil,
		StoragePath:              nil,
		ConcurrencyLock: ConcurrencyLock{
			Enabled: false,
		},
		AllowedMethods: nil,
		UseStale:       nil,
		Durations:      nil,
	}
}
