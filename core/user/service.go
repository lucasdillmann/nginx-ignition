package user

import (
	"context"

	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/common/configuration"
	"dillmann.com.br/nginx-ignition/core/common/coreerror"
	"dillmann.com.br/nginx-ignition/core/common/pagination"
	"dillmann.com.br/nginx-ignition/core/common/validation"
	"dillmann.com.br/nginx-ignition/core/user/passwordhash"
)

var invalidCredentialsError = coreerror.New("Invalid username or password", true)

type service struct {
	repository    Repository
	configuration *configuration.Configuration
}

func (s *service) authenticate(ctx context.Context, username string, password string) (*User, error) {
	usr, err := s.repository.FindByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	if usr == nil {
		return nil, invalidCredentialsError
	}

	passwordMatches, err := passwordhash.New(s.configuration).Verify(password, usr.PasswordHash, usr.PasswordSalt)
	if err != nil {
		return nil, err
	}

	if !passwordMatches {
		return nil, invalidCredentialsError
	}

	return usr, nil
}

func (s *service) changePassword(ctx context.Context, id uuid.UUID, currentPassword string, newPassword string) error {
	hash := passwordhash.New(s.configuration)
	databaseState, err := s.repository.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if databaseState == nil {
		return coreerror.New("No user found with provided ID", true)
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
		return validation.SingleFieldError("currentPassword", "Not your current password")
	}

	if len(newPassword) < minimumPasswordLength {
		return validation.SingleFieldError("newPassword", "Must have at least 8 characters")
	}

	updatedHash, updatedSalt, err := hash.Hash(newPassword)
	if err != nil {
		return err
	}

	databaseState.PasswordHash = updatedHash
	databaseState.PasswordSalt = updatedSalt
	return s.repository.Save(ctx, databaseState)
}

func (s *service) getById(ctx context.Context, id uuid.UUID) (*User, error) {
	return s.repository.FindByID(ctx, id)
}

func (s *service) deleteById(ctx context.Context, id uuid.UUID) error {
	return s.repository.DeleteByID(ctx, id)
}

func (s *service) count(ctx context.Context) (int, error) {
	return s.repository.Count(ctx)
}

func (s *service) isOnboardingCompleted(ctx context.Context) (bool, error) {
	count, err := s.repository.Count(ctx)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (s *service) save(ctx context.Context, request *SaveRequest, currentUserId *uuid.UUID) error {
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

	if err := newValidator(s.repository).validate(ctx, updatedState, databaseState, request, currentUserId); err != nil {
		return err
	}

	return s.repository.Save(ctx, updatedState)
}

func (s *service) isEnabled(ctx context.Context, id uuid.UUID) (bool, error) {
	return s.repository.IsEnabledByID(ctx, id)
}

func (s *service) list(ctx context.Context, pageSize, pageNumber int, searchTerms *string) (*pagination.Page[User], error) {
	return s.repository.FindPage(ctx, pageSize, pageNumber, searchTerms)
}

func (s *service) resetPassword(ctx context.Context, username string) (string, error) {
	user, err := s.repository.FindByUsername(ctx, username)
	if err != nil {
		return "", err
	}

	if user == nil {
		return "", coreerror.New("User not found", true)
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
