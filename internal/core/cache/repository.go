package cache

import (
	"context"

	"github.com/google/uuid"

	"nginx-ignition/internal/core/common/pagination"
)

type Repository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*Cache, error)
	InUseByID(ctx context.Context, id uuid.UUID) (bool, error)
	ExistsByID(ctx context.Context, id uuid.UUID) (bool, error)
	DeleteByID(ctx context.Context, id uuid.UUID) error
	FindPage(
		ctx context.Context,
		pageNumber, pageSize int,
		searchTerms *string,
	) (*pagination.Page[Cache], error)
	FindAllInUse(ctx context.Context) ([]Cache, error)
	Save(ctx context.Context, cache *Cache) error
}
