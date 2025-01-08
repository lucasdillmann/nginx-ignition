package database

import (
	"database/sql"
	"dillmann.com.br/nginx-ignition/core/common/configuration"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

type Database struct {
	configuration *configuration.Configuration
	db            *sql.DB
}

func New(configuration configuration.Configuration) *Database {
	return &Database{&configuration, nil}
}

func (d *Database) BeginTransaction() (*sql.Tx, error) {
	return (*d.db).Begin()
}

func (d *Database) close() {
	log.Println("Closing the database connection")

	if d.db == nil {
		return
	}

	if err := d.db.Close(); err != nil {
		log.Printf("Unable to close the connection to database: %s", err)
	}
}

func (d *Database) init() error {
	cfg := (*d.configuration).WithPrefix("nginx-ignition.database")

	var host, port, username, password, driver, sslMode, name string
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

	if driver, err = cfg.Get("driver"); err != nil {
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

	log.Printf(
		"Starting database connection to %s on %s:%s using username %s, driver %s and SSL mode %s",
		name, host, port, username, driver, sslMode,
	)

	connectionParams := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s sslmode=%s dbname=%s",
		host, port, username, password, sslMode, name,
	)

	if d.db, err = sql.Open(driver, connectionParams); err != nil {
		return err
	}

	if err = d.db.Ping(); err != nil {
		return err
	}

	return nil
}
