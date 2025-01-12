package migrations

import (
	"dillmann.com.br/nginx-ignition/core/common/configuration"
	"dillmann.com.br/nginx-ignition/database/common/database"
)

const (
	tableName = "schema_version"
)

type migrations struct {
	database      *database.Database
	configuration *configuration.Configuration
}

func newMigrations(database *database.Database, configuration *configuration.Configuration) *migrations {
	return &migrations{
		database:      database,
		configuration: configuration.WithPrefix("nginx-ignition.database"),
	}
}

func (m *migrations) migrate() error {
	driverName, err := m.configuration.Get("driver")
	if err != nil {
		return err
	}

	db := m.database.Unwrap()

	if err := m.migrateFromFlyway(db, driverName); err != nil {
		return err
	}

	return m.runScripts(db, driverName)
}
