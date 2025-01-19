package certificate

import (
	"dillmann.com.br/nginx-ignition/core/common/pagination"
	"github.com/google/uuid"
)

type DeleteCommand func(id uuid.UUID) error

type AvailableProvidersCommand func() ([]*AvailableProvider, error)

type GetCommand func(id uuid.UUID) (*Certificate, error)

type ListCommand func(pageSize, pageNumber int, searchTerms *string) (*pagination.Page[*Certificate], error)

type IssueCommand func(request *IssueRequest) (*Certificate, error)

type RenewCommand func(id uuid.UUID) error
