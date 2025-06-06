package certificate

import (
	"context"
	"dillmann.com.br/nginx-ignition/core/common/pagination"
	"github.com/google/uuid"
)

type Commands struct {
	Delete             func(ctx context.Context, id uuid.UUID) error
	AvailableProviders func(ctx context.Context) ([]*AvailableProvider, error)
	Get                func(ctx context.Context, id uuid.UUID) (*Certificate, error)
	List               func(ctx context.Context, pageSize, pageNumber int, searchTerms *string) (*pagination.Page[*Certificate], error)
	Issue              func(ctx context.Context, request *IssueRequest) (*Certificate, error)
	Renew              func(ctx context.Context, id uuid.UUID) error
}
