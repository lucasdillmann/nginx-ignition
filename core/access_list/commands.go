package access_list

import (
	"dillmann.com.br/nginx-ignition/core/common/pagination"
	"github.com/google/uuid"
)

type DeleteCommand func(id uuid.UUID) error

type GetCommand func(id uuid.UUID) (*AccessList, error)

type ListCommand func(pageSize, pageNumber int, searchTerms *string) (*pagination.Page[AccessList], error)

type SaveCommand func(accessList *AccessList) error
