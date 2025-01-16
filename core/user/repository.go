package user

import (
	"dillmann.com.br/nginx-ignition/core/common/pagination"
	"github.com/google/uuid"
)

type Repository interface {
	Save(user *User) error
	DeleteByID(id uuid.UUID) error
	FindByID(id uuid.UUID) (*User, error)
	FindByUsername(username string) (*User, error)
	FindPage(pageSize, pageNumber int, searchTerms *string) (*pagination.Page[*User], error)
	IsEnabledByID(id uuid.UUID) (bool, error)
	Count() (int, error)
}
