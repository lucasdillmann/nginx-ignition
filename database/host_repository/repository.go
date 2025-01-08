package host_repository

import (
	"dillmann.com.br/nginx-ignition/core/common/core_errors"
	"dillmann.com.br/nginx-ignition/core/host"
	"dillmann.com.br/nginx-ignition/database/common/database"
	"github.com/google/uuid"
)

type repository struct {
	database *database.Database
}

func New(database *database.Database) host.Repository {
	return &repository{database}
}

func (r repository) ExistsByAccessListId(_ uuid.UUID) (bool, error) {
	return false, core_errors.NotImplemented()
}
