package core

import (
	"github.com/lucasdillmann/nginx-ignition/internal/core/accesslist"
	"github.com/lucasdillmann/nginx-ignition/internal/core/backup"
	"github.com/lucasdillmann/nginx-ignition/internal/core/binding"
	"github.com/lucasdillmann/nginx-ignition/internal/core/cache"
	"github.com/lucasdillmann/nginx-ignition/internal/core/certificate"
	"github.com/lucasdillmann/nginx-ignition/internal/core/common/broadcast"
	"github.com/lucasdillmann/nginx-ignition/internal/core/common/container"
	"github.com/lucasdillmann/nginx-ignition/internal/core/common/scheduler"
	"github.com/lucasdillmann/nginx-ignition/internal/core/host"
	"github.com/lucasdillmann/nginx-ignition/internal/core/integration"
	"github.com/lucasdillmann/nginx-ignition/internal/core/nginx"
	"github.com/lucasdillmann/nginx-ignition/internal/core/settings"
	"github.com/lucasdillmann/nginx-ignition/internal/core/stream"
	"github.com/lucasdillmann/nginx-ignition/internal/core/user"
	"github.com/lucasdillmann/nginx-ignition/internal/core/vpn"
)

func Install() error {
	return container.Run(
		broadcast.Install,
		scheduler.Install,
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
