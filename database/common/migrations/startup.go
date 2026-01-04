package migrations

import (
	"context"

	"dillmann.com.br/nginx-ignition/core/common/lifecycle"
)

type startup struct {
	migrations *Migrations
}

func registerStartup(lc *lifecycle.Lifecycle, mig *Migrations) {
	command := &startup{mig}
	lc.RegisterStartup(command)
}

func (d startup) Priority() int {
	return startupPriority
}

func (d startup) Async() bool {
	return false
}

func (d startup) Run(_ context.Context) error {
	return d.migrations.Migrate()
}
