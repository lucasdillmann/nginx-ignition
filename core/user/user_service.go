package user

import (
	"dillmann.com.br/nginx-ignition/core/common/configuration"
	"dillmann.com.br/nginx-ignition/core/common/core_errors"
	"dillmann.com.br/nginx-ignition/core/common/pagination"
	"dillmann.com.br/nginx-ignition/core/user/password_hash"
	"errors"
	"github.com/google/uuid"
)

type service struct {
	repository    *Repository
	configuration *configuration.Configuration
}

func (s *service) authenticate(_ string, _ string) (*User, error) {
	return nil, core_errors.NotImplemented()
}

func (s *service) changePassword(id uuid.UUID, currentPassword string, newPassword string) error {
	hash := password_hash.New(s.configuration)
	databaseState, err := (*s.repository).FindById(id)
	if err != nil {
		return err
	}

	if databaseState == nil {
		return errors.New("no user exists with the provided ID")
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
		return errors.New("current password does not match")
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
	return (*s.repository).FindById(id)
}

func (s *service) deleteById(id uuid.UUID) error {
	return (*s.repository).DeleteById(id)
}

func (s *service) count() (int64, error) {
	return (*s.repository).Count()
}

func (s *service) save(request *SaveRequest, currentUserId uuid.UUID) error {
	var passwordHash, passwordSalt string
	var databaseState *User
	var err error

	if request.Id != uuid.Nil && request.Password == "" {
		databaseState, err = (*s.repository).FindById(request.Id)
		if err != nil {
			return err
		}

		if databaseState != nil {
			passwordHash = databaseState.PasswordHash
			passwordSalt = databaseState.PasswordSalt
		}
	} else if request.Password != "" {
		passwordHash, passwordSalt, err = password_hash.New(s.configuration).Hash(request.Password)
		if err != nil {
			return err
		}
	}

	updatedState := &User{
		Id:           request.Id,
		Enabled:      request.Enabled,
		Name:         request.Name,
		Username:     request.Username,
		PasswordHash: passwordHash,
		PasswordSalt: passwordSalt,
		Role:         request.Role,
	}

	if err := newValidator().validate(s.repository, updatedState, databaseState, request, currentUserId); err != nil {
		return err
	}

	return (*s.repository).Save(updatedState)
}

func (s *service) isEnabled(id uuid.UUID) (bool, error) {
	return (*s.repository).IsEnabledById(id)
}

func (s *service) list(pageSize int64, pageNumber int64, searchTerms string) (*pagination.Page[User], error) {
	return (*s.repository).FindPage(pageSize, pageNumber, searchTerms)
}
