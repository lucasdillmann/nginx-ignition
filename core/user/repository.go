package user

import (
	"dillmann.com.br/nginx-ignition/core/common/pagination"
	"github.com/google/uuid"
)

type Repository interface {
	Save(user *User) error
	DeleteById(id uuid.UUID) error
	FindById(id uuid.UUID) (*User, error)
	FindByUsername(username string) (*User, error)
	FindPage(pageSize int64, pageNumber int64, searchTerms string) (*pagination.Page[User], error)
	IsEnabledById(id uuid.UUID) (bool, error)
	Count() (int64, error)
}
