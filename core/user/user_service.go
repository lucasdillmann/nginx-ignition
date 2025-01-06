package user

import (
	"dillmann.com.br/nginx-ignition/core/common/configuration"
	"dillmann.com.br/nginx-ignition/core/common/core_errors"
	"dillmann.com.br/nginx-ignition/core/common/pagination"
	"github.com/google/uuid"
)

type service struct {
	repository    *Repository
	configuration *configuration.Configuration
}

func (s *service) authenticate(_ string, _ string) (*User, error) {
	return nil, core_errors.NotImplemented()
}

func (s *service) changePassword(id uuid.UUID, _ string, _ string) error {
	_, err := (*s.repository).FindById(id)
	if err != nil {
		return err
	}

	return core_errors.NotImplemented()
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

func (s *service) save(user *User) error {
	if err := newValidator().validate(user); err != nil {
		return err
	}

	return (*s.repository).Save(user)
}

func (s *service) isEnabled(id uuid.UUID) (bool, error) {
	return (*s.repository).IsEnabledById(id)
}

func (s *service) list(pageSize int64, pageNumber int64, searchTerms string) (*pagination.Page[User], error) {
	return (*s.repository).FindPage(pageSize, pageNumber, searchTerms)
}
