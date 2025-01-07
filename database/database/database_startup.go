package database

import (
	"dillmann.com.br/nginx-ignition/core/common/lifecycle"
)

type startup struct {
	database *Database
}

func Register(lifecycle *lifecycle.Lifecycle, database *Database) {
	command := &startup{database}
	lifecycle.RegisterStartup(command)
}

func (d startup) Priority() int {
	return 0
}

func (d startup) Async() bool {
	return false
}

func (d startup) Run() error {
	return d.database.init()
}
