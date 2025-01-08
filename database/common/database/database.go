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

	var host, port, username, password, driver, sslmode string
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

	if sslmode, err = cfg.Get("ssl-mode"); err != nil || sslmode != "disable" {
		sslmode = "require"
	}

	port, _ = cfg.Get("port")

	log.Printf(
		"Starting database connection to %s using username %s, driver %s and SSL mode %s",
		host, username, driver, sslmode,
	)

	connectionParams := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s sslmode=%s",
		host, port, username, password, sslmode,
	)

	if d.db, err = sql.Open(driver, connectionParams); err != nil {
		return err
	}

	if err = d.db.Ping(); err != nil {
		return err
	}

	return nil
}
