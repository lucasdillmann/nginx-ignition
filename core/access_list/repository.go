package access_list

import (
	"dillmann.com.br/nginx-ignition/core/common/pagination"
	"github.com/google/uuid"
)

type Repository interface {
	FindByID(id uuid.UUID) (*AccessList, error)
	DeleteByID(id uuid.UUID) error
	FindByName(name string) (*AccessList, error)
	FindPage(pageNumber, pageSize int, searchTerms *string) (*pagination.Page[AccessList], error)
	FindAll() ([]*AccessList, error)
	Save(accessList *AccessList) error
}
