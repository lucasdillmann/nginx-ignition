package user

import (
	"context"

	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/common/configuration"
	"dillmann.com.br/nginx-ignition/core/common/coreerror"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
	"dillmann.com.br/nginx-ignition/core/common/pagination"
	"dillmann.com.br/nginx-ignition/core/common/validation"
	"dillmann.com.br/nginx-ignition/core/user/passwordhash"
)

type service struct {
	repository    Repository
	configuration *configuration.Configuration
}

func newService(repository Repository, cfg *configuration.Configuration) *service {
	return &service{
		repository:    repository,
		configuration: cfg,
	}
}

func (s *service) Authenticate(ctx context.Context, username, password string) (*User, error) {
	usr, err := s.repository.FindByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	if usr == nil {
		return nil, coreerror.New(
			i18n.M(ctx, i18n.K.UserErrorInvalidCredentials),
			true,
		)
	}

	passwordMatches, err := passwordhash.New(s.configuration).
		Verify(password, usr.PasswordHash, usr.PasswordSalt)
	if err != nil {
		return nil, err
	}

	if !passwordMatches {
		return nil, coreerror.New(
			i18n.M(ctx, i18n.K.UserErrorInvalidCredentials),
			true,
		)
	}

	return usr, nil
}

func (s *service) UpdatePassword(
	ctx context.Context,
	id uuid.UUID,
	currentPassword, newPassword string,
) error {
	hash := passwordhash.New(s.configuration)
	databaseState, err := s.repository.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if databaseState == nil {
		return coreerror.New(i18n.M(ctx, i18n.K.UserErrorNotFoundByID), true)
	}

	passwordMatches, err := hash.Verify(
		currentPassword,
		databaseState.PasswordHash,
		databaseState.PasswordSalt,
	)
	if err != nil {
		return err
	}

	if !passwordMatches {
		return validation.SingleFieldError(
			"currentPassword",
			i18n.M(ctx, i18n.K.UserValidationCurrentPasswordMismatch),
		)
	}

	if len(newPassword) < minimumPasswordLength {
		return validation.SingleFieldError(
			"newPassword",
			i18n.M(ctx, i18n.K.CommonValidationTooShort).V("min", minimumPasswordLength),
		)
	}

	updatedHash, updatedSalt, err := hash.Hash(newPassword)
	if err != nil {
		return err
	}

	databaseState.PasswordHash = updatedHash
	databaseState.PasswordSalt = updatedSalt
	return s.repository.Save(ctx, databaseState)
}

func (s *service) Get(ctx context.Context, id uuid.UUID) (*User, error) {
	return s.repository.FindByID(ctx, id)
}

func (s *service) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repository.DeleteByID(ctx, id)
}

func (s *service) GetCount(ctx context.Context) (int, error) {
	return s.repository.Count(ctx)
}

func (s *service) OnboardingCompleted(ctx context.Context) (bool, error) {
	count, err := s.repository.Count(ctx)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (s *service) Save(ctx context.Context, request *SaveRequest, currentUserID *uuid.UUID) error {
	var passwordHash, passwordSalt string
	var databaseState *User
	var err error

	databaseState, err = s.repository.FindByID(ctx, request.ID)
	if err != nil {
		return err
	}

	if request.Password == nil && databaseState != nil {
		passwordHash = databaseState.PasswordHash
		passwordSalt = databaseState.PasswordSalt
	} else if request.Password != nil {
		passwordHash, passwordSalt, err = passwordhash.New(s.configuration).Hash(*request.Password)
		if err != nil {
			return err
		}
	}

	updatedState := &User{
		ID:           request.ID,
		Enabled:      request.Enabled,
		Name:         request.Name,
		Username:     request.Username,
		PasswordHash: passwordHash,
		PasswordSalt: passwordSalt,
		Permissions:  request.Permissions,
	}

	if err := newValidator(s.repository).validate(
		ctx,
		updatedState,
		databaseState,
		request,
		currentUserID,
	); err != nil {
		return err
	}

	return s.repository.Save(ctx, updatedState)
}

func (s *service) GetStatus(ctx context.Context, id uuid.UUID) (bool, error) {
	return s.repository.IsEnabledByID(ctx, id)
}

func (s *service) List(
	ctx context.Context,
	pageSize, pageNumber int,
	searchTerms *string,
) (*pagination.Page[User], error) {
	return s.repository.FindPage(ctx, pageSize, pageNumber, searchTerms)
}

func (s *service) resetPassword(ctx context.Context, username string) (string, error) {
	user, err := s.repository.FindByUsername(ctx, username)
	if err != nil {
		return "", err
	}

	if user == nil {
		return "", coreerror.New(i18n.M(ctx, i18n.K.UserErrorNotFound), true)
	}

	newPassword := uuid.NewString()[:8]
	updatedHash, updatedSalt, err := passwordhash.New(s.configuration).Hash(newPassword)
	if err != nil {
		return "", err
	}

	user.PasswordHash = updatedHash
	user.PasswordSalt = updatedSalt
	return newPassword, s.repository.Save(ctx, user)
}
