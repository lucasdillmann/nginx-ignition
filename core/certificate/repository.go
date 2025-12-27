package certificate

import (
	"context"

	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/common/pagination"
)

type Repository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*Certificate, error)
	ExistsByID(ctx context.Context, id uuid.UUID) (bool, error)
	InUseByID(ctx context.Context, id uuid.UUID) (bool, error)
	DeleteByID(ctx context.Context, id uuid.UUID) error
	Save(ctx context.Context, certificate *Certificate) error
	FindPage(ctx context.Context, pageSize, pageNumber int, searchTerms *string) (*pagination.Page[Certificate], error)
	FindAllDueToRenew(ctx context.Context) ([]Certificate, error)
	GetAutoRenewSettings(ctx context.Context) (*AutoRenewSettings, error)
}
