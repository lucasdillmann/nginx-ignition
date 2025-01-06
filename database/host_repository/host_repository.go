package host_repository

import (
	"database/sql"
	"dillmann.com.br/nginx-ignition/core/common/core_errors"
	"dillmann.com.br/nginx-ignition/core/host"
	"github.com/google/uuid"
)

type repository struct {
	database *sql.DB
}

func New(database *sql.DB) host.Repository {
	return &repository{database}
}

func (r repository) ExistsByAccessListId(_ uuid.UUID) (bool, error) {
	return false, core_errors.NotImplemented()
}
