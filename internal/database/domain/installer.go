package domain

import (
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/container"
	"github.com/lucasdillmann/nginx-ignition/internal/database/domain/accesslist"
	"github.com/lucasdillmann/nginx-ignition/internal/database/domain/backup"
	"github.com/lucasdillmann/nginx-ignition/internal/database/domain/cache"
	"github.com/lucasdillmann/nginx-ignition/internal/database/domain/certificate"
	"github.com/lucasdillmann/nginx-ignition/internal/database/domain/host"
	"github.com/lucasdillmann/nginx-ignition/internal/database/domain/integration"
	"github.com/lucasdillmann/nginx-ignition/internal/database/domain/settings"
	"github.com/lucasdillmann/nginx-ignition/internal/database/domain/stream"
	"github.com/lucasdillmann/nginx-ignition/internal/database/domain/user"
	"github.com/lucasdillmann/nginx-ignition/internal/database/domain/vpn"
)

func Install() error {
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
