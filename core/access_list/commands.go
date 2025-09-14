package access_list

import (
	"context"

	"dillmann.com.br/nginx-ignition/core/common/pagination"
	"github.com/google/uuid"
)

type Commands struct {
	Delete func(ctx context.Context, id uuid.UUID) error
	Get    func(ctx context.Context, id uuid.UUID) (*AccessList, error)
	List   func(ctx context.Context, pageSize, pageNumber int, searchTerms *string) (*pagination.Page[*AccessList], error)
	Save   func(ctx context.Context, accessList *AccessList) error
}
