package certificate

import (
	"context"

	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/common/pagination"
)

type Commands interface {
	Delete(ctx context.Context, id uuid.UUID) error
	AvailableProviders(ctx context.Context) ([]AvailableProvider, error)
	Exists(ctx context.Context, id uuid.UUID) (bool, error)
	Get(ctx context.Context, id uuid.UUID) (*Certificate, error)
	List(
		ctx context.Context,
		pageSize, pageNumber int,
		searchTerms *string,
	) (*pagination.Page[Certificate], error)
	Issue(ctx context.Context, request *IssueRequest) (*Certificate, error)
	Renew(ctx context.Context, id uuid.UUID) error
}
