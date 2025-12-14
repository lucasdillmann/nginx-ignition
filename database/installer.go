package database

import (
	"dillmann.com.br/nginx-ignition/core/common/container"
	"dillmann.com.br/nginx-ignition/database/accesslist"
	"dillmann.com.br/nginx-ignition/database/backup"
	"dillmann.com.br/nginx-ignition/database/certificate/client"
	"dillmann.com.br/nginx-ignition/database/certificate/server"
	"dillmann.com.br/nginx-ignition/database/common/database"
	"dillmann.com.br/nginx-ignition/database/common/migrations"
	"dillmann.com.br/nginx-ignition/database/host"
	"dillmann.com.br/nginx-ignition/database/integration"
	"dillmann.com.br/nginx-ignition/database/settings"
	"dillmann.com.br/nginx-ignition/database/stream"
	"dillmann.com.br/nginx-ignition/database/user"
	"dillmann.com.br/nginx-ignition/database/vpn"
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
		host.New,
		user.New,
		settings.New,
		server.New,
		client.New,
		integration.New,
		stream.New,
		backup.New,
		vpn.New,
	)
}
