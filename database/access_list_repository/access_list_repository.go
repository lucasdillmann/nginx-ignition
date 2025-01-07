package access_list_repository

import (
	"dillmann.com.br/nginx-ignition/core/access_list"
	"dillmann.com.br/nginx-ignition/core/common/core_errors"
	"dillmann.com.br/nginx-ignition/core/common/pagination"
	"dillmann.com.br/nginx-ignition/database/database"
	"github.com/google/uuid"
)

type repository struct {
	database *database.Database
}

func New(database *database.Database) access_list.Repository {
	return &repository{database}
}

func (r repository) FindById(_ uuid.UUID) (*access_list.AccessList, error) {
	return nil, core_errors.NotImplemented()
}

func (r repository) DeleteById(_ uuid.UUID) error {
	return core_errors.NotImplemented()
}

func (r repository) FindByName(_ string) (*access_list.AccessList, error) {
	return nil, core_errors.NotImplemented()
}

func (r repository) FindPage(
	_ int64,
	_ int64,
	_ string,
) (*pagination.Page[access_list.AccessList], error) {
	return nil, core_errors.NotImplemented()
}

func (r repository) FindAll() (*[]access_list.AccessList, error) {
	return nil, core_errors.NotImplemented()
}

func (r repository) Save(_ *access_list.AccessList) error {
	return core_errors.NotImplemented()
}
