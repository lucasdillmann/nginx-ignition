package database

import (
	"database/sql"
	"dillmann.com.br/nginx-ignition/core/common/configuration"
	"dillmann.com.br/nginx-ignition/core/common/core_errors"
	"log"
)

type Database struct {
	configuration *configuration.Configuration
	db            *sql.DB
}

func New(configuration configuration.Configuration) *Database {
	return &Database{&configuration, nil}
}

func (d *Database) init() error {
	log.Println("Starting database connection")
	return core_errors.NotImplemented()
}
