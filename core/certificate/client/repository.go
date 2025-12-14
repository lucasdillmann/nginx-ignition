package client

import (
	"context"

	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/common/pagination"
)

type Repository interface {
	IsInUseByID(ctx context.Context, id uuid.UUID) (bool, error)
	FindByID(ctx context.Context, id uuid.UUID) (*Certificate, error)
	DeleteByID(ctx context.Context, id uuid.UUID) error
	Save(ctx context.Context, certificate *Certificate) error
	FindPage(ctx context.Context, pageSize, pageNumber int, searchTerms *string) (*pagination.Page[*Certificate], error)
}
