package user

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/lucasdillmann/nginx-ignition/internal/business/core/pagination"
)

type Commands interface {
	Authenticate(
		ctx context.Context,
		username, password, code string,
	) (AuthenticationOutcome, *User, error)
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
	FinishOnboarding(ctx context.Context, request *SaveRequest) error
	UpdatePassword(
		ctx context.Context,
		id, currentSessionID uuid.UUID,
		oldPassword, newPassword string,
	) error
	UpdateProfile(ctx context.Context, id uuid.UUID, name, username string) error
	OnboardingCompleted(ctx context.Context) (bool, error)
	GetTOTPStatus(ctx context.Context, id uuid.UUID) (bool, error)
	DisableTOTP(ctx context.Context, id uuid.UUID) error
	EnableTOTP(ctx context.Context, id uuid.UUID) (string, error)
	ActivateTOTP(ctx context.Context, id uuid.UUID, code string) (bool, error)
	ListAPITokens(
		ctx context.Context,
		id uuid.UUID,
		pageSize, pageNumber int,
		searchTerms *string,
	) (*pagination.Page[APIToken], error)
	CreateAPIToken(
		ctx context.Context,
		id uuid.UUID,
		request *NewAPITokenRequest,
	) (*APIToken, error)
	DeleteAPIToken(ctx context.Context, id, tokenID uuid.UUID) error
	FindAPIToken(ctx context.Context, id, tokenID uuid.UUID) (*APIToken, error)
	StartSession(
		ctx context.Context,
		userID uuid.UUID,
		expiresAt time.Time,
	) (*Session, error)
	FindSession(ctx context.Context, id uuid.UUID) (*Session, error)
	StopSession(ctx context.Context, id uuid.UUID) error
	RotateSession(
		ctx context.Context,
		oldID, userID uuid.UUID,
		expiresAt time.Time,
	) (*Session, error)
}
