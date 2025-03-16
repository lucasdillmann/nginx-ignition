package user

import (
	"context"
	"dillmann.com.br/nginx-ignition/core/common/pagination"
	"github.com/google/uuid"
)

type Repository interface {
	Save(ctx context.Context, user *User) error
	DeleteByID(ctx context.Context, id uuid.UUID) error
	FindByID(ctx context.Context, id uuid.UUID) (*User, error)
	FindByUsername(ctx context.Context, username string) (*User, error)
	FindPage(ctx context.Context, pageSize, pageNumber int, searchTerms *string) (*pagination.Page[*User], error)
	IsEnabledByID(ctx context.Context, id uuid.UUID) (bool, error)
	Count(ctx context.Context) (int, error)
}
