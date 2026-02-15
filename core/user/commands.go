package user

import (
	"context"

	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/common/pagination"
)

type Commands interface {
	Authenticate(ctx context.Context, username, password string) (*User, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Get(ctx context.Context, id uuid.UUID) (*User, error)
	GetCount(ctx context.Context) (int, error)
	GetStatus(ctx context.Context, id uuid.UUID) (bool, error)
	List(
		ctx context.Context,
		pageSize, pageNumber int,
		searchTerms *string,
	) (*pagination.Page[User], error)
	Save(ctx context.Context, user *SaveRequest, currentUserID *uuid.UUID) error
	UpdatePassword(ctx context.Context, id uuid.UUID, oldPassword, newPassword string) error
	OnboardingCompleted(ctx context.Context) (bool, error)
	GetTOTPStatus(ctx context.Context, id uuid.UUID) (bool, error)
	DisableTOTP(ctx context.Context, id uuid.UUID) error
	EnableTOTP(ctx context.Context, id uuid.UUID) (string, error)
	ActivateTOTP(ctx context.Context, id uuid.UUID, code string) (bool, error)
}
