//nolint:revive
package api

import (
	"github.com/lucasdillmann/nginx-ignition/internal/api/accesslist"
	"github.com/lucasdillmann/nginx-ignition/internal/api/backup"
	"github.com/lucasdillmann/nginx-ignition/internal/api/cache"
	"github.com/lucasdillmann/nginx-ignition/internal/api/certificate"
	"github.com/lucasdillmann/nginx-ignition/internal/api/common/server"
	"github.com/lucasdillmann/nginx-ignition/internal/api/frontend"
	"github.com/lucasdillmann/nginx-ignition/internal/api/healthcheck"
	"github.com/lucasdillmann/nginx-ignition/internal/api/host"
	"github.com/lucasdillmann/nginx-ignition/internal/api/i18n"
	"github.com/lucasdillmann/nginx-ignition/internal/api/integration"
	"github.com/lucasdillmann/nginx-ignition/internal/api/nginx"
	"github.com/lucasdillmann/nginx-ignition/internal/api/settings"
	"github.com/lucasdillmann/nginx-ignition/internal/api/stream"
	"github.com/lucasdillmann/nginx-ignition/internal/api/user"
	"github.com/lucasdillmann/nginx-ignition/internal/api/vpn"
	"github.com/lucasdillmann/nginx-ignition/internal/core/common/container"
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
