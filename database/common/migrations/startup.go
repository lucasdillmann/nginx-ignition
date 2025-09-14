package migrations

import (
	"context"

	"dillmann.com.br/nginx-ignition/core/common/lifecycle"
)

type startup struct {
	migrations *migrations
}

func registerStartup(lifecycle *lifecycle.Lifecycle, migrations *migrations) {
	command := &startup{migrations}
	lifecycle.RegisterStartup(command)
}

func (d startup) Priority() int {
	return startupPriority
}

func (d startup) Async() bool {
	return false
}

func (d startup) Run(_ context.Context) error {
	return d.migrations.migrate()
}
