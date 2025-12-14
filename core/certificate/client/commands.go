package client

import (
	"context"

	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/common/pagination"
)

type Commands struct {
	List         func(ctx context.Context, pageSize, pageNumber int, searchTerms *string) (*pagination.Page[*Certificate], error)
	Create       func(ctx context.Context, request *CreateRequest) (*Certificate, error)
	Delete       func(ctx context.Context, id uuid.UUID) error
	Get          func(ctx context.Context, id uuid.UUID) (*Certificate, error)
	Update       func(ctx context.Context, id uuid.UUID, request *UpdateRequest) error
	ReplaceCA    func(ctx context.Context, id uuid.UUID, request *ReplaceCARequest) error
	UpdateClient func(ctx context.Context, id uuid.UUID, request *UpdateClientRequest) error
	CreateClient func(ctx context.Context, id uuid.UUID, request *CreateClientRequest) error
	DeleteClient func(ctx context.Context, id uuid.UUID, clientId uuid.UUID) error
}
