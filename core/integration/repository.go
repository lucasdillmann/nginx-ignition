package integration

import (
	"context"

	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/common/pagination"
)

type Repository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*Integration, error)
	ExistsByName(ctx context.Context, name string) (*bool, error)
	ExistsByID(ctx context.Context, id uuid.UUID) (*bool, error)
	InUseByID(ctx context.Context, id uuid.UUID) (*bool, error)
	Save(ctx context.Context, integration *Integration) error
	DeleteByID(ctx context.Context, id uuid.UUID) error
	FindPage(
		ctx context.Context,
		pageSize, pageNumber int,
		searchTerms *string,
		enabledOnly bool,
	) (*pagination.Page[*Integration], error)
}
