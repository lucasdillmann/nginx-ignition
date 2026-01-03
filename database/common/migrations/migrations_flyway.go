package migrations

import (
	"database/sql"

	"dillmann.com.br/nginx-ignition/core/common/log"
)

func (m *Migrations) migrateFromFlyway(db *sql.DB, driver string) error {
	if driver != "postgres" {
		return nil
	}

	result, err := db.Query(`
		SELECT column_name 
		FROM information_schema.columns
		WHERE table_name = 'schema_version'
		AND table_schema = 'public'
		AND column_name = 'installed_rank'
	`)
	if err != nil {
		return err
	}

	//nolint:errcheck
	defer result.Close()

	if result.Next() {
		return m.convertFlywayTable(db)
	}

	return nil
}

func (m *Migrations) convertFlywayTable(db *sql.DB) error {
	var err error

	_, err = db.Exec("ALTER TABLE schema_version RENAME TO schema_version_backup")
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		CREATE TABLE schema_version AS
		SELECT version, false AS dirty
		FROM schema_version_backup
		ORDER BY installed_rank DESC
		LIMIT 1
	`)
	if err != nil {
		return err
	}

	log.Infof("Database migration history successfully converted from Flyway to Go Migrate")
	return nil
}
