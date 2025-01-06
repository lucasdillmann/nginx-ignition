package user_repository

import (
	"database/sql"
	"dillmann.com.br/nginx-ignition/core/common/core_errors"
	"dillmann.com.br/nginx-ignition/core/common/pagination"
	"dillmann.com.br/nginx-ignition/core/user"
	"github.com/google/uuid"
)

type repository struct {
	database *sql.DB
}

func New(database *sql.DB) user.Repository {
	return &repository{database}
}

func (r repository) Save(_ *user.User) error {
	return core_errors.NotImplemented()
}

func (r repository) DeleteById(_ uuid.UUID) error {
	return core_errors.NotImplemented()
}

func (r repository) FindById(_ uuid.UUID) (*user.User, error) {
	return nil, core_errors.NotImplemented()
}

func (r repository) FindByUsername(_ string) (*user.User, error) {
	return nil, core_errors.NotImplemented()
}

func (r repository) FindPage(_ int64, _ int64, _ string) (*pagination.Page[user.User], error) {
	return nil, core_errors.NotImplemented()
}

func (r repository) IsEnabledById(_ uuid.UUID) (bool, error) {
	return false, core_errors.NotImplemented()
}

func (r repository) Count() (int64, error) {
	return 0, core_errors.NotImplemented()
}
