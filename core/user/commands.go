package user

import (
	"dillmann.com.br/nginx-ignition/core/common/pagination"
	"github.com/google/uuid"
)

type AuthenticateCommand func(username string, password string) (*User, error)

type DeleteCommand func(id uuid.UUID) error

type GetCommand func(id uuid.UUID) (*User, error)

type GetCountCommand func() (int, error)

type GetStatusCommand func(id uuid.UUID) (bool, error)

type ListCommand func(pageSize, pageNumber int, searchTerms *string) (*pagination.Page[User], error)

type SaveCommand func(user *SaveRequest, currentUserId *uuid.UUID) error

type UpdatePasswordCommand func(id uuid.UUID, oldPassword string, newPassword string) error

type OnboardingCompletedCommand func() (bool, error)
