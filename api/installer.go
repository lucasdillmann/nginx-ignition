package api

import (
	"dillmann.com.br/nginx-ignition/api/accesslist"
	"dillmann.com.br/nginx-ignition/api/backup"
	"dillmann.com.br/nginx-ignition/api/cache"
	"dillmann.com.br/nginx-ignition/api/certificate"
	"dillmann.com.br/nginx-ignition/api/common/server"
	"dillmann.com.br/nginx-ignition/api/frontend"
	"dillmann.com.br/nginx-ignition/api/healthcheck"
	"dillmann.com.br/nginx-ignition/api/host"
	"dillmann.com.br/nginx-ignition/api/integration"
	"dillmann.com.br/nginx-ignition/api/nginx"
	"dillmann.com.br/nginx-ignition/api/settings"
	"dillmann.com.br/nginx-ignition/api/stream"
	"dillmann.com.br/nginx-ignition/api/user"
	"dillmann.com.br/nginx-ignition/api/vpn"
	"dillmann.com.br/nginx-ignition/core/common/container"
)

func Install() error {
	return container.Run(
		server.Install,
		healthcheck.Install,
		settings.Install,
		accesslist.Install,
		certificate.Install,
		user.Install,
		host.Install,
		integration.Install,
		nginx.Install,
		stream.Install,
		backup.Install,
		vpn.Install,
		cache.Install,
		frontend.Install,
	)
}
