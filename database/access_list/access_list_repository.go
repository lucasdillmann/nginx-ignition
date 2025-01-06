package access_list_repository

import (
	"database/sql"
	"dillmann.com.br/nginx-ignition/core/access_list"
	"dillmann.com.br/nginx-ignition/core/common/pagination"
	"errors"
	"github.com/google/uuid"
)

type accessListDatabaseRepository struct {
	database *sql.DB
}

func New(database *sql.DB) access_list.AccessListRepository {
	return &accessListDatabaseRepository{
		database: database,
	}
}

func (repository accessListDatabaseRepository) FindById(_ uuid.UUID) (*access_list.AccessList, error) {
	return nil, notImplemented()
}

func (repository accessListDatabaseRepository) DeleteById(_ uuid.UUID) error {
	return notImplemented()
}

func (repository accessListDatabaseRepository) FindByName(_ string) (*access_list.AccessList, error) {
	return nil, notImplemented()
}

func (repository accessListDatabaseRepository) FindPage(
	_ int64,
	_ int64,
	_ string,
) (*pagination.Page[access_list.AccessList], error) {
	return nil, notImplemented()
}

func (repository accessListDatabaseRepository) FindAll() (*[]access_list.AccessList, error) {
	return nil, notImplemented()
}

func (repository accessListDatabaseRepository) Save(_ *access_list.AccessList) error {
	return notImplemented()
}

func notImplemented() error {
	return errors.New("not yet implemented")
}
