package stream

import (
	"context"

	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/common/pagination"
)

type Commands struct {
	Save          func(ctx context.Context, input *Stream) error
	Delete        func(ctx context.Context, id uuid.UUID) error
	List          func(ctx context.Context, pageSize, pageNumber int, searchTerms *string) (*pagination.Page[*Stream], error)
	Get           func(ctx context.Context, id uuid.UUID) (*Stream, error)
	GetAllEnabled func(ctx context.Context) ([]*Stream, error)
	Exists        func(ctx context.Context, id uuid.UUID) (bool, error)
}
