package accesslist

import (
	"context"

	"github.com/google/uuid"

	"nginx-ignition/internal/core/common/pagination"
)

type Repository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*AccessList, error)
	ExistsByID(ctx context.Context, id uuid.UUID) (bool, error)
	InUseByID(ctx context.Context, id uuid.UUID) (bool, error)
	DeleteByID(ctx context.Context, id uuid.UUID) error
	FindPage(
		ctx context.Context,
		pageNumber, pageSize int,
		searchTerms *string,
	) (*pagination.Page[AccessList], error)
	FindAll(ctx context.Context) ([]AccessList, error)
	Save(ctx context.Context, accessList *AccessList) error
}
