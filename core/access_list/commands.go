package access_list

import (
	"context"
	"dillmann.com.br/nginx-ignition/core/common/pagination"
	"github.com/google/uuid"
)

type DeleteCommand func(ctx context.Context, id uuid.UUID) error

type GetCommand func(ctx context.Context, id uuid.UUID) (*AccessList, error)

type ListCommand func(ctx context.Context, pageSize, pageNumber int, searchTerms *string) (*pagination.Page[*AccessList], error)

type SaveCommand func(ctx context.Context, accessList *AccessList) error
