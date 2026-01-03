package stream

import (
	"context"

	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/common/pagination"
)

type Repository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*Stream, error)
	DeleteByID(ctx context.Context, id uuid.UUID) error
	Save(ctx context.Context, stream *Stream) error
	FindPage(
		ctx context.Context,
		pageSize, pageNumber int,
		searchTerms *string,
	) (*pagination.Page[Stream], error)
	FindAllEnabled(ctx context.Context) ([]Stream, error)
	ExistsByID(ctx context.Context, id uuid.UUID) (bool, error)
}
