package access_list

import (
	"dillmann.com.br/nginx-ignition/core/common/pagination"
	"github.com/google/uuid"
)

type DeleteByIdCommand func(id uuid.UUID) error

type GetByIdCommand func(id uuid.UUID) (*AccessList, error)

type ListCommand func(pageSize int64, pageNumber int64, searchTerms string) (*pagination.Page[AccessList], error)

type SaveCommand func(accessList *AccessList) error
