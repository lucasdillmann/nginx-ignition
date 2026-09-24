package integration

import (
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/container"
	"github.com/lucasdillmann/nginx-ignition/internal/integration/domain"
)

func Install() error {
	return container.Run(
		domain.Install,
	)
}
