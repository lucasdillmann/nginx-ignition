package cache

import (
	"context"

	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/common/pagination"
)

type Commands struct {
	Delete func(ctx context.Context, id uuid.UUID) error
	Get    func(ctx context.Context, id uuid.UUID) (*Cache, error)
	List   func(ctx context.Context, pageSize, pageNumber int, searchTerms *string) (*pagination.Page[*Cache], error)
	Save   func(ctx context.Context, cache *Cache) error
}
