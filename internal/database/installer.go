package database

import (
	"nginx-ignition/internal/core/common/container"
	"nginx-ignition/internal/database/accesslist"
	"nginx-ignition/internal/database/backup"
	"nginx-ignition/internal/database/cache"
	"nginx-ignition/internal/database/certificate"
	"nginx-ignition/internal/database/common/database"
	"nginx-ignition/internal/database/common/migrations"
	"nginx-ignition/internal/database/host"
	"nginx-ignition/internal/database/integration"
	"nginx-ignition/internal/database/settings"
	"nginx-ignition/internal/database/stream"
	"nginx-ignition/internal/database/user"
	"nginx-ignition/internal/database/vpn"
)

func Install() error {
	if err := container.Run(
		database.Install,
		migrations.Install,
	); err != nil {
		return err
	}

	return container.Provide(
		accesslist.New,
		cache.New,
		host.New,
		user.New,
		settings.New,
		certificate.New,
		integration.New,
		stream.New,
		backup.New,
		vpn.New,
	)
}
