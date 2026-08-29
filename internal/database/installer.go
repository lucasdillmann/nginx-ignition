package database

import (
	"github.com/lucasdillmann/nginx-ignition/internal/core/common/container"
	"github.com/lucasdillmann/nginx-ignition/internal/database/accesslist"
	"github.com/lucasdillmann/nginx-ignition/internal/database/backup"
	"github.com/lucasdillmann/nginx-ignition/internal/database/cache"
	"github.com/lucasdillmann/nginx-ignition/internal/database/certificate"
	"github.com/lucasdillmann/nginx-ignition/internal/database/common/database"
	"github.com/lucasdillmann/nginx-ignition/internal/database/common/migrations"
	"github.com/lucasdillmann/nginx-ignition/internal/database/host"
	"github.com/lucasdillmann/nginx-ignition/internal/database/integration"
	"github.com/lucasdillmann/nginx-ignition/internal/database/settings"
	"github.com/lucasdillmann/nginx-ignition/internal/database/stream"
	"github.com/lucasdillmann/nginx-ignition/internal/database/user"
	"github.com/lucasdillmann/nginx-ignition/internal/database/vpn"
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
