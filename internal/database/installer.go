package database

import (
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/container"
	"github.com/lucasdillmann/nginx-ignition/internal/database/core"
	"github.com/lucasdillmann/nginx-ignition/internal/database/domain"
)

func Install() error {
	return container.Run(
		core.Install,
		domain.Install,
	)
}
