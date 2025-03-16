package user

import (
	"context"
	"dillmann.com.br/nginx-ignition/core/common/pagination"
	"github.com/google/uuid"
)

type AuthenticateCommand func(ctx context.Context, username string, password string) (*User, error)

type DeleteCommand func(ctx context.Context, id uuid.UUID) error

type GetCommand func(ctx context.Context, id uuid.UUID) (*User, error)

type GetCountCommand func(ctx context.Context) (int, error)

type GetStatusCommand func(ctx context.Context, id uuid.UUID) (bool, error)

type ListCommand func(ctx context.Context, pageSize, pageNumber int, searchTerms *string) (*pagination.Page[*User], error)

type SaveCommand func(ctx context.Context, user *SaveRequest, currentUserId *uuid.UUID) error

type UpdatePasswordCommand func(ctx context.Context, id uuid.UUID, oldPassword string, newPassword string) error

type OnboardingCompletedCommand func(ctx context.Context) (bool, error)
