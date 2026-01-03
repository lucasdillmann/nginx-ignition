package stream

import (
	"context"

	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/common/pagination"
)

type Commands interface {
	Save(ctx context.Context, input *Stream) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(
		ctx context.Context,
		pageSize, pageNumber int,
		searchTerms *string,
	) (*pagination.Page[Stream], error)
	Get(ctx context.Context, id uuid.UUID) (*Stream, error)
	GetAllEnabled(ctx context.Context) ([]Stream, error)
	Exists(ctx context.Context, id uuid.UUID) (bool, error)
}
