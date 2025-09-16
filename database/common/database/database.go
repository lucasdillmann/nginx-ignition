package database

import (
	"database/sql"

	"github.com/uptrace/bun"

	"dillmann.com.br/nginx-ignition/core/common/configuration"
	"dillmann.com.br/nginx-ignition/core/common/log"
)

type Database struct {
	configuration    *configuration.Configuration
	db               *sql.DB
	bun              *bun.DB
	connectionString string
}

func newDatabase(configuration *configuration.Configuration) *Database {
	return &Database{
		configuration: configuration.WithPrefix("nginx-ignition.database"),
		db:            nil,
		bun:           nil,
	}
}

func (d *Database) Unwrap() *sql.DB {
	return d.db
}

func (d *Database) Begin() (bun.Tx, error) {
	return d.bun.Begin()
}

func (d *Database) Select() *bun.SelectQuery {
	return d.bun.NewSelect()
}

func (d *Database) ConnectionString() string {
	return d.connectionString
}

func (d *Database) close() {
	log.Infof("Closing the database connection")

	if d.db == nil {
		return
	}

	if err := d.db.Close(); err != nil {
		log.Warnf("Unable to close the connection to database: %s", err)
	}
}
