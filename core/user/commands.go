package user

import (
	"context"

	"dillmann.com.br/nginx-ignition/core/common/pagination"
	"github.com/google/uuid"
)

type Commands struct {
	Authenticate        func(ctx context.Context, username string, password string) (*User, error)
	Delete              func(ctx context.Context, id uuid.UUID) error
	Get                 func(ctx context.Context, id uuid.UUID) (*User, error)
	GetCount            func(ctx context.Context) (int, error)
	GetStatus           func(ctx context.Context, id uuid.UUID) (bool, error)
	List                func(ctx context.Context, pageSize, pageNumber int, searchTerms *string) (*pagination.Page[*User], error)
	Save                func(ctx context.Context, user *SaveRequest, currentUserId *uuid.UUID) error
	UpdatePassword      func(ctx context.Context, id uuid.UUID, oldPassword string, newPassword string) error
	OnboardingCompleted func(ctx context.Context) (bool, error)
}
