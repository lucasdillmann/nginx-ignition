package database

import (
	"dillmann.com.br/nginx-ignition/core/common/lifecycle"
)

type shutdown struct {
	database *Database
}

func registerShutdown(lifecycle *lifecycle.Lifecycle, database *Database) {
	command := &shutdown{database}
	lifecycle.RegisterShutdown(command)
}

func (d shutdown) Priority() int {
	return shutdownPriority
}

func (d shutdown) Run() {
	d.database.close()
}
