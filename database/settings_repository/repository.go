package settings_repository

import (
	"dillmann.com.br/nginx-ignition/core/common/core_errors"
	"dillmann.com.br/nginx-ignition/core/settings"
	"dillmann.com.br/nginx-ignition/database/common/database"
)

type repository struct {
	database *database.Database
}

func New(database *database.Database) settings.Repository {
	return &repository{database}
}

func (r repository) Save(_ *settings.Settings) error {
	return core_errors.NotImplemented()
}
