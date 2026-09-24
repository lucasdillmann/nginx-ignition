package core

import (
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/container"
	"github.com/lucasdillmann/nginx-ignition/internal/database/core/database"
	"github.com/lucasdillmann/nginx-ignition/internal/database/core/migrations"
)

func Install() error {
	return container.Run(
		database.Install,
		migrations.Install,
	)
}
