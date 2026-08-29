package core

import (
	"nginx-ignition/internal/core/accesslist"
	"nginx-ignition/internal/core/backup"
	"nginx-ignition/internal/core/binding"
	"nginx-ignition/internal/core/cache"
	"nginx-ignition/internal/core/certificate"
	"nginx-ignition/internal/core/common/broadcast"
	"nginx-ignition/internal/core/common/container"
	"nginx-ignition/internal/core/common/scheduler"
	"nginx-ignition/internal/core/host"
	"nginx-ignition/internal/core/integration"
	"nginx-ignition/internal/core/nginx"
	"nginx-ignition/internal/core/settings"
	"nginx-ignition/internal/core/stream"
	"nginx-ignition/internal/core/user"
	"nginx-ignition/internal/core/vpn"
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
