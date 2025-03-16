package certificate

import (
	"context"
	"dillmann.com.br/nginx-ignition/core/common/pagination"
	"github.com/google/uuid"
)

type DeleteCommand func(ctx context.Context, id uuid.UUID) error

type AvailableProvidersCommand func(ctx context.Context) ([]*AvailableProvider, error)

type GetCommand func(ctx context.Context, id uuid.UUID) (*Certificate, error)

type ListCommand func(ctx context.Context, pageSize, pageNumber int, searchTerms *string) (*pagination.Page[*Certificate], error)

type IssueCommand func(ctx context.Context, request *IssueRequest) (*Certificate, error)

type RenewCommand func(ctx context.Context, id uuid.UUID) error
