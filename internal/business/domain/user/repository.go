package user

import (
	"context"

	"github.com/google/uuid"

	"github.com/lucasdillmann/nginx-ignition/internal/business/core/pagination"
)

type Repository interface {
	Save(ctx context.Context, user *User) error
	DeleteByID(ctx context.Context, id uuid.UUID) error
	FindByID(ctx context.Context, id uuid.UUID) (*User, error)
	FindByUsername(ctx context.Context, username string) (*User, error)
	FindPage(
		ctx context.Context,
		pageSize, pageNumber int,
		searchTerms *string,
	) (*pagination.Page[User], error)
	CreateToken(ctx context.Context, token *APIToken) error
	DeleteTokenByID(ctx context.Context, userID, id uuid.UUID) error
	ExistsTokenByName(ctx context.Context, userID uuid.UUID, name string) (bool, error)
	FindTokenByID(ctx context.Context, userID, id uuid.UUID) (*APIToken, error)
	FindTokensByUserID(
		ctx context.Context,
		userID uuid.UUID,
		pageNumber, pageSize int,
		searchTerms *string,
	) (*pagination.Page[APIToken], error)
	IsEnabledByID(ctx context.Context, id uuid.UUID) (bool, error)
	Count(ctx context.Context) (int, error)
	TryCreateInitialUser(ctx context.Context, user *User) (bool, error)
	TryUpdateLastUsedTOTPCode(ctx context.Context, id uuid.UUID, code string) (bool, error)
}
