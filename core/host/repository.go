package host

import (
	"context"

	"dillmann.com.br/nginx-ignition/core/common/pagination"
	"github.com/google/uuid"
)

type Repository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*Host, error)
	DeleteByID(ctx context.Context, id uuid.UUID) error
	Save(ctx context.Context, host *Host) error
	FindPage(ctx context.Context, pageSize, pageNumber int, searchTerms *string) (*pagination.Page[*Host], error)
	FindAllEnabled(ctx context.Context) ([]*Host, error)
	FindDefault(ctx context.Context) (*Host, error)
	ExistsByID(ctx context.Context, id uuid.UUID) (bool, error)
	ExistsByCertificateID(ctx context.Context, certificateId uuid.UUID) (bool, error)
	ExistsCertificateByID(ctx context.Context, certificateId uuid.UUID) (bool, error)
	ExistsByAccessListID(ctx context.Context, accessListId uuid.UUID) (bool, error)
}
