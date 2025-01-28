package database

import (
	"dillmann.com.br/nginx-ignition/core/common/lifecycle"
)

type startup struct {
	database *Database
}

func registerStartup(lifecycle *lifecycle.Lifecycle, database *Database) {
	lifecycle.RegisterStartup(startup{database})
}

func (d startup) Priority() int {
	return startupPriority
}

func (d startup) Async() bool {
	return false
}

func (d startup) Run() error {
	return d.database.init()
}
