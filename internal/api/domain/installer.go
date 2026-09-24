package domain

import (
	"github.com/lucasdillmann/nginx-ignition/internal/api/domain/accesslist"
	"github.com/lucasdillmann/nginx-ignition/internal/api/domain/backup"
	"github.com/lucasdillmann/nginx-ignition/internal/api/domain/cache"
	"github.com/lucasdillmann/nginx-ignition/internal/api/domain/certificate"
	"github.com/lucasdillmann/nginx-ignition/internal/api/domain/frontend"
	"github.com/lucasdillmann/nginx-ignition/internal/api/domain/healthcheck"
	"github.com/lucasdillmann/nginx-ignition/internal/api/domain/host"
	"github.com/lucasdillmann/nginx-ignition/internal/api/domain/i18n"
	"github.com/lucasdillmann/nginx-ignition/internal/api/domain/integration"
	"github.com/lucasdillmann/nginx-ignition/internal/api/domain/nginx"
	"github.com/lucasdillmann/nginx-ignition/internal/api/domain/settings"
	"github.com/lucasdillmann/nginx-ignition/internal/api/domain/stream"
	"github.com/lucasdillmann/nginx-ignition/internal/api/domain/user"
	"github.com/lucasdillmann/nginx-ignition/internal/api/domain/vpn"
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/container"
)

func Install() error {
	return container.Run(
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
