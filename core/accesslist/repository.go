package accesslist

import (
	"context"

	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/common/pagination"
)

type Repository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*AccessList, error)
	DeleteByID(ctx context.Context, id uuid.UUID) error
	FindByName(ctx context.Context, name string) (*AccessList, error)
	FindPage(ctx context.Context, pageNumber, pageSize int, searchTerms *string) (*pagination.Page[*AccessList], error)
	FindAll(ctx context.Context) ([]*AccessList, error)
	Save(ctx context.Context, accessList *AccessList) error
	IsInUseByID(ctx context.Context, id uuid.UUID) (bool, error)
}
