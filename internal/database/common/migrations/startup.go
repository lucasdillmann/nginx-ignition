package migrations

import (
	"context"

	"github.com/lucasdillmann/nginx-ignition/internal/core/common/lifecycle"
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
