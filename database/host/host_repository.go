package host_repository

import (
	"database/sql"
	"dillmann.com.br/nginx-ignition/core/host"
	"errors"
	"github.com/google/uuid"
)

type hostDatabaseRepository struct {
	database *sql.DB
}

func New(database *sql.DB) host.HostRepository {
	return &hostDatabaseRepository{
		database: database,
	}
}

func (repository hostDatabaseRepository) ExistsByAccessListId(_ uuid.UUID) (bool, error) {
	return false, errors.New("not yet implemented")
}
