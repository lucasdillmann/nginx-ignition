package database

import (
	"context"

	"dillmann.com.br/nginx-ignition/core/common/lifecycle"
)

type shutdown struct {
	database *Database
}

func registerShutdown(lc *lifecycle.Lifecycle, db *Database) {
	lc.RegisterShutdown(shutdown{db})
}

func (d shutdown) Priority() int {
	return shutdownPriority
}

func (d shutdown) Run(_ context.Context) {
	d.database.close()
}
