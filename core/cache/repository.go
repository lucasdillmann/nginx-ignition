package cache

import (
	"context"

	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/common/pagination"
)

type Repository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*Cache, error)
	IsInUseByID(ctx context.Context, id uuid.UUID) (bool, error)
	DeleteByID(ctx context.Context, id uuid.UUID) error
	FindByName(ctx context.Context, name string) (*Cache, error)
	FindPage(ctx context.Context, pageNumber, pageSize int, searchTerms *string) (*pagination.Page[*Cache], error)
	FindAll(ctx context.Context) ([]*Cache, error)
	Save(ctx context.Context, cache *Cache) error
}
