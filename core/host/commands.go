package host

import (
	"dillmann.com.br/nginx-ignition/core/common/pagination"
	"dillmann.com.br/nginx-ignition/core/common/validation"
	"github.com/google/uuid"
)

type SaveCommand = func(input *Host) error
type DeleteCommand = func(id uuid.UUID) error
type ListCommand = func(pageSize, pageNumber int, searchTerms *string) (*pagination.Page[*Host], error)
type GetCommand = func(id uuid.UUID) (*Host, error)
type GetAllEnabledCommand = func() ([]*Host, error)
type ExistsCommand func(id uuid.UUID) (bool, error)
type ValidateBindingCommand func(
	path string,
	index int,
	binding *Binding,
	context *validation.ConsistencyValidator,
) error
