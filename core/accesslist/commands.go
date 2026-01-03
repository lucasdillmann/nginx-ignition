package accesslist

import (
	"context"

	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/common/pagination"
)

type Commands interface {
	Delete(ctx context.Context, id uuid.UUID) error
	Get(ctx context.Context, id uuid.UUID) (*AccessList, error)
	GetAll(ctx context.Context) ([]AccessList, error)
	Exists(ctx context.Context, id uuid.UUID) (bool, error)
	List(
		ctx context.Context,
		pageSize, pageNumber int,
		searchTerms *string,
	) (*pagination.Page[AccessList], error)
	Save(ctx context.Context, accessList *AccessList) error
}
