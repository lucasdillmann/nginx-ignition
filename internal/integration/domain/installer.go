package domain

import (
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/container"
	"github.com/lucasdillmann/nginx-ignition/internal/integration/domain/docker"
	"github.com/lucasdillmann/nginx-ignition/internal/integration/domain/truenas"
)

func Install() error {
	return container.Run(
		docker.Install,
		truenas.Install,
	)
}
