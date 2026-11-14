package migrations

import (
	"database/sql"
	"errors"
	"fmt"
	"path"

	"github.com/lucasdillmann/nginx-ignition/internal/business/core/log"

	"github.com/golang-migrate/migrate/v4"
	migratordb "github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

func (m *Migrations) runScripts(db *sql.DB, driverName string) error {
	var err error
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

	scriptsPath := path.Join("scripts", driverName)
	log.Infof("Running embedded database migrations from %s", scriptsPath)

	sourceInstance, err := iofs.New(scriptsFS, scriptsPath)
	if err != nil {
		return err
	}

	migrator, err := migrate.NewWithInstance("iofs", sourceInstance, driverName, driverInstance)
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
