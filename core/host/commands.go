package host

import (
	"context"

	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/common/pagination"
)

type Commands interface {
	Save(ctx context.Context, input *Host) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(
		ctx context.Context,
		pageSize, pageNumber int,
		searchTerms *string,
	) (*pagination.Page[Host], error)
	Get(ctx context.Context, id uuid.UUID) (*Host, error)
	GetAllEnabled(ctx context.Context) ([]Host, error)
	Exists(ctx context.Context, id uuid.UUID) (bool, error)
}
