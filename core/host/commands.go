package host

import (
	"context"
	"dillmann.com.br/nginx-ignition/core/common/pagination"
	"dillmann.com.br/nginx-ignition/core/common/validation"
	"github.com/google/uuid"
)

type SaveCommand = func(ctx context.Context, input *Host) error
type DeleteCommand = func(ctx context.Context, id uuid.UUID) error
type ListCommand = func(ctx context.Context, pageSize, pageNumber int, searchTerms *string) (*pagination.Page[*Host], error)
type GetCommand = func(ctx context.Context, id uuid.UUID) (*Host, error)
type GetAllEnabledCommand = func(ctx context.Context) ([]*Host, error)
type ExistsCommand func(ctx context.Context, id uuid.UUID) (bool, error)
type ValidateBindingCommand func(
	ctx context.Context,
	path string,
	index int,
	binding *Binding,
	context *validation.ConsistencyValidator,
) error
