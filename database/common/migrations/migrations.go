package migrations

import (
	"dillmann.com.br/nginx-ignition/core/common/log"
	"dillmann.com.br/nginx-ignition/database/common/database"
)

type migrations struct {
	database *database.Database
}

func newMigrations(database *database.Database) *migrations {
	return &migrations{database}
}

func (m *migrations) migrate() error {
	log.Warn("Database migrations aren't implemented yet. You need to run the SQL script manually for now.")
	return nil
}
