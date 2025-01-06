package access_list

import (
	"dillmann.com.br/nginx-ignition/core/common/pagination"
	"github.com/google/uuid"
)

type AccessListRepository interface {
	FindById(id uuid.UUID) (*AccessList, error)
	DeleteById(id uuid.UUID) error
	FindByName(name string) (*AccessList, error)
	FindPage(pageNumber int64, pageSize int64, searchTerms string) (*pagination.Page[AccessList], error)
	FindAll() (*[]AccessList, error)
	Save(accessList *AccessList) error
}
