package host

import (
	"context"

	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/common/pagination"
)

type Commands struct {
	Save          func(ctx context.Context, input *Host) error
	Delete        func(ctx context.Context, id uuid.UUID) error
	List          func(ctx context.Context, pageSize, pageNumber int, searchTerms *string) (*pagination.Page[Host], error)
	Get           func(ctx context.Context, id uuid.UUID) (*Host, error)
	GetAllEnabled func(ctx context.Context) ([]Host, error)
	Exists        func(ctx context.Context, id uuid.UUID) (bool, error)
}
