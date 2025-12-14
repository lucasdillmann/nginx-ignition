package core

import (
	"dillmann.com.br/nginx-ignition/core/accesslist"
	"dillmann.com.br/nginx-ignition/core/backup"
	"dillmann.com.br/nginx-ignition/core/certificate/client"
	"dillmann.com.br/nginx-ignition/core/certificate/server"
	"dillmann.com.br/nginx-ignition/core/common/broadcast"
	"dillmann.com.br/nginx-ignition/core/common/container"
	"dillmann.com.br/nginx-ignition/core/common/scheduler"
	"dillmann.com.br/nginx-ignition/core/host"
	"dillmann.com.br/nginx-ignition/core/integration"
	"dillmann.com.br/nginx-ignition/core/nginx"
	"dillmann.com.br/nginx-ignition/core/settings"
	"dillmann.com.br/nginx-ignition/core/stream"
	"dillmann.com.br/nginx-ignition/core/user"
	"dillmann.com.br/nginx-ignition/core/vpn"
)

func Install() error {
	return container.Run(
		broadcast.Install,
		scheduler.Install,
		settings.Install,
		user.Install,
		accesslist.Install,
		server.Install,
		client.Install,
		vpn.Install,
		host.Install,
		integration.Install,
		stream.Install,
		nginx.Install,
		backup.Install,
	)
}
