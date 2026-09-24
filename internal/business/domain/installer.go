package domain

import (
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/container"
	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/accesslist"
	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/backup"
	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/binding"
	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/cache"
	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/certificate"
	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/host"
	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/integration"
	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/nginx"
	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/settings"
	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/stream"
	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/user"
	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/vpn"
)

func Install() error {
	return container.Run(
		settings.Install,
		user.Install,
		accesslist.Install,
		binding.Install,
		cache.Install,
		certificate.Install,
		vpn.Install,
		host.Install,
		integration.Install,
		stream.Install,
		nginx.Install,
		backup.Install,
	)
}
