package migrations

import (
	"database/sql"

	"dillmann.com.br/nginx-ignition/core/common/configuration"
	"dillmann.com.br/nginx-ignition/core/common/log"
	"dillmann.com.br/nginx-ignition/database/common/database"
)

const (
	tableName = "schema_version"
)

type Migrations struct {
	database      *database.Database
	configuration *configuration.Configuration
}

func New(db *database.Database, cfg *configuration.Configuration) *Migrations {
	return &Migrations{
		database:      db,
		configuration: cfg.WithPrefix("nginx-ignition.database"),
	}
}

func (m *Migrations) Migrate() error {
	driverName, err := m.configuration.Get("driver")
	if err != nil {
		return err
	}

	db := m.database.Unwrap()

	if err := m.migrateFromFlyway(db, driverName); err != nil {
		return err
	}

	m.unsetDirtyStatus(db)
	return m.runScripts(db, driverName)
}

func (m *Migrations) unsetDirtyStatus(db *sql.DB) {
	rows, err := db.Query("select version, dirty from schema_version limit 1")
	if err != nil {
		return
	}

	//nolint:errcheck
	defer rows.Close()

	if rows.Next() {
		var version int
		var dirty bool

		if err = rows.Scan(&version, &dirty); err != nil {
			return
		}

		if dirty {
			_, err = db.Exec("update schema_version set dirty = false, version = $1", version-1)
			if err != nil {
				log.Warnf("Unable to unset database dirty status: %s", err.Error())
				return
			}
		}
	}
}
