package access_list

import (
	"dillmann.com.br/nginx-ignition/core/common/pagination"
	"github.com/google/uuid"
)

type DeleteAccessListByIdCommand func(id uuid.UUID) error

type GetAccessListByIdCommand func(id uuid.UUID) (*AccessList, error)

type ListAccessListCommand func(
	pageSize int64,
	pageNumber int64,
	searchTerms string,
) (*pagination.Page[AccessList], error)

type SaveAccessListCommand func(accessList *AccessList) error
