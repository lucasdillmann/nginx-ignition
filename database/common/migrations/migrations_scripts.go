package migrations

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	migratordb "github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"

	"dillmann.com.br/nginx-ignition/core/common/log"

	// Loads the file driver to be dynamically used by gomigrate
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func (m *migrations) runScripts(db *sql.DB, driverName string) error {
	scriptsPath, err := m.configuration.Get("migrations-path")
	if err != nil {
		return err
	}

	var driverInstance migratordb.Driver
	switch driverName {
	case "postgres":
		driverInstance, err = postgres.WithInstance(db, &postgres.Config{
			MigrationsTable: tableName,
		})
	case "sqlite":
		driverInstance, err = sqlite3.WithInstance(db, &sqlite3.Config{
			MigrationsTable: tableName,
		})
	default:
		return fmt.Errorf("unsupported driver: %s", driverName)
	}

	if err != nil {
		return err
	}

	fullPath := fmt.Sprintf("%s/%s", scriptsPath, driverName)
	log.Infof("Running database migrations from %s", fullPath)

	migrator, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", fullPath),
		driverName,
		driverInstance,
	)
	if err != nil {
		return err
	}

	err = migrator.Up()
	if errors.Is(err, migrate.ErrNoChange) {
		log.Infof("Database is up-to-date. No changes needed.")
		return nil
	}

	if err == nil {
		installedVersion, _, _ := migrator.Version()
		log.Infof("Database updated to version %d", installedVersion)
	}

	return err
}
