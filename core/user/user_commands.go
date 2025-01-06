package user

import (
	"dillmann.com.br/nginx-ignition/core/common/pagination"
	"github.com/google/uuid"
)

type AuthenticateCommand func(username string, password string) (*User, error)

type DeleteByIdCommand func(id uuid.UUID) error

type GetByIdCommand func(id uuid.UUID) (*User, error)

type GetCountCommand func() (int64, error)

type GetStatusCommand func(id uuid.UUID) (bool, error)

type ListCommand func(pageSize int64, pageNumber int64, searchTerms string) (*pagination.Page[User], error)

type SaveCommand func(user *User) error

type ChangePasswordCommand func(id uuid.UUID, oldPassword string, newPassword string) error
