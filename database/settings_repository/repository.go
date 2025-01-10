package settings_repository

import (
	"dillmann.com.br/nginx-ignition/core/settings"
	"dillmann.com.br/nginx-ignition/database/common/database"
)

type repository struct {
	database *database.Database
}

func New(database *database.Database) settings.Repository {
	return &repository{database}
}
