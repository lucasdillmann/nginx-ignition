package cache

import (
	"context"

	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/common/pagination"
)

type Commands interface {
	Delete(ctx context.Context, id uuid.UUID) error
	Get(ctx context.Context, id uuid.UUID) (*Cache, error)
	Exists(ctx context.Context, id uuid.UUID) (bool, error)
	List(
		ctx context.Context,
		pageSize, pageNumber int,
		searchTerms *string,
	) (*pagination.Page[Cache], error)
	GetAllInUse(ctx context.Context) ([]Cache, error)
	Save(ctx context.Context, cache *Cache) error
}
