package database

import (
	"database/sql"
	"dillmann.com.br/nginx-ignition/core/common/core_error"
	"dillmann.com.br/nginx-ignition/core/common/log"
	"fmt"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/schema"
	"os"

	_ "github.com/lib/pq"
	_ "modernc.org/sqlite"
)

func (d *Database) init() error {
	var driver string
	var err error

	if driver, err = d.configuration.Get("driver"); err != nil {
		return err
	}

	switch driver {
	case "postgres":
		return d.initPostgres()
	case "sqlite":
		return d.initSqlite()
	default:
		return core_error.New(fmt.Sprintf("unsupported database driver: %s", driver), true)
	}
}

func (d *Database) initPostgres() error {
	cfg := *d.configuration

	var host, port, username, password, sslMode, name string
	var err error

	if host, err = cfg.Get("host"); err != nil {
		return err
	}

	if username, err = cfg.Get("username"); err != nil {
		return err
	}

	if password, err = cfg.Get("password"); err != nil {
		return err
	}

	if name, err = cfg.Get("name"); err != nil {
		return err
	}

	if sslMode, err = cfg.Get("ssl-mode"); err != nil || sslMode != "disable" {
		sslMode = "require"
	}

	if port, err = cfg.Get("port"); err != nil {
		return err
	}

	log.Infof(
		"Starting database connection to %s on %s:%s using username %s and SSL mode %s",
		name, host, port, username, sslMode,
	)

	connectionParams := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s sslmode=%s dbname=%s",
		host, port, username, password, sslMode, name,
	)

	return d.initBun("postgres", connectionParams, pgdialect.New())
}

func (d *Database) initSqlite() error {
	log.Warnf(
		"Application is configured to use the embedded SQLite database. This isn't recommended for " +
			"production environments, please refer to the documentation in order to migrate to PostgreSQL.",
	)

	folder, err := d.configuration.Get("data-path")
	if err != nil {
		return err
	}

	filePath := fmt.Sprintf("%s/nginx-ignition.db", folder)
	log.Infof("Starting database connection to SQLite on %s", filePath)

	if err = os.MkdirAll(folder, os.ModePerm); err != nil {
		return err
	}

	return d.initBun(
		"sqlite",
		fmt.Sprintf("file:%s", filePath),
		sqlitedialect.New(),
	)
}

func (d *Database) initBun(driverName, connectionParams string, dialect schema.Dialect) error {
	var err error

	if d.db, err = sql.Open(driverName, connectionParams); err != nil {
		return err
	}

	if err = d.db.Ping(); err != nil {
		return err
	}

	d.bun = bun.NewDB(d.db, dialect)
	return nil
}
