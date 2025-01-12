package user

import (
	"dillmann.com.br/nginx-ignition/core/common/configuration"
	"dillmann.com.br/nginx-ignition/core/common/core_error"
	"dillmann.com.br/nginx-ignition/core/common/pagination"
	"dillmann.com.br/nginx-ignition/core/common/validation"
	"dillmann.com.br/nginx-ignition/core/user/password_hash"
	"github.com/google/uuid"
)

type service struct {
	repository    *Repository
	configuration *configuration.Configuration
}

func (s *service) authenticate(username string, password string) (*User, error) {
	usr, err := (*s.repository).FindByUsername(username)
	if err != nil {
		return nil, err
	}

	passwordMatches, err := password_hash.New(s.configuration).Verify(password, usr.PasswordHash, usr.PasswordSalt)
	if err != nil {
		return nil, err
	}

	if !passwordMatches {
		return nil, core_error.New("Invalid username or password", true)
	}

	return usr, nil
}

func (s *service) changePassword(id uuid.UUID, currentPassword string, newPassword string) error {
	hash := password_hash.New(s.configuration)
	databaseState, err := (*s.repository).FindByID(id)
	if err != nil {
		return err
	}

	if databaseState == nil {
		return core_error.New("No user found with provided ID", true)
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
	return (*s.repository).Save(databaseState)
}

func (s *service) getById(id uuid.UUID) (*User, error) {
	return (*s.repository).FindByID(id)
}

func (s *service) deleteById(id uuid.UUID) error {
	return (*s.repository).DeleteByID(id)
}

func (s *service) count() (int, error) {
	return (*s.repository).Count()
}

func (s *service) isOnboardingCompleted() (bool, error) {
	count, err := (*s.repository).Count()
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (s *service) save(request *SaveRequest, currentUserId *uuid.UUID) error {
	var passwordHash, passwordSalt string
	var databaseState *User
	var err error

	databaseState, err = (*s.repository).FindByID(request.ID)
	if err != nil {
		return err
	}

	if request.Password == nil && databaseState != nil {
		passwordHash = databaseState.PasswordHash
		passwordSalt = databaseState.PasswordSalt
	} else if request.Password != nil {
		passwordHash, passwordSalt, err = password_hash.New(s.configuration).Hash(*request.Password)

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
		Role:         request.Role,
	}

	if err := newValidator(s.repository).validate(updatedState, databaseState, request, currentUserId); err != nil {
		return err
	}

	return (*s.repository).Save(updatedState)
}

func (s *service) isEnabled(id uuid.UUID) (bool, error) {
	return (*s.repository).IsEnabledByID(id)
}

func (s *service) list(pageSize, pageNumber int, searchTerms *string) (*pagination.Page[User], error) {
	return (*s.repository).FindPage(pageSize, pageNumber, searchTerms)
}
