package database

import (
	"context"

	"dillmann.com.br/nginx-ignition/core/common/lifecycle"
)

type startup struct {
	database *Database
}

func registerStartup(lc *lifecycle.Lifecycle, db *Database) {
	lc.RegisterStartup(startup{db})
}

func (d startup) Priority() int {
	return startupPriority
}

func (d startup) Async() bool {
	return false
}

func (d startup) Run(ctx context.Context) error {
	return d.database.Init(ctx)
}
