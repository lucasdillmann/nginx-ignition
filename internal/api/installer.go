//nolint:revive
package api

import (
	"nginx-ignition/internal/api/accesslist"
	"nginx-ignition/internal/api/backup"
	"nginx-ignition/internal/api/cache"
	"nginx-ignition/internal/api/certificate"
	"nginx-ignition/internal/api/common/server"
	"nginx-ignition/internal/api/frontend"
	"nginx-ignition/internal/api/healthcheck"
	"nginx-ignition/internal/api/host"
	"nginx-ignition/internal/api/i18n"
	"nginx-ignition/internal/api/integration"
	"nginx-ignition/internal/api/nginx"
	"nginx-ignition/internal/api/settings"
	"nginx-ignition/internal/api/stream"
	"nginx-ignition/internal/api/user"
	"nginx-ignition/internal/api/vpn"
	"nginx-ignition/internal/core/common/container"
)

func Install() error {
	return container.Run(
		server.Install,
		healthcheck.Install,
		settings.Install,
		accesslist.Install,
		cache.Install,
		certificate.Install,
		user.Install,
		host.Install,
		i18n.Install,
		integration.Install,
		nginx.Install,
		stream.Install,
		backup.Install,
		vpn.Install,
		frontend.Install,
	)
}
