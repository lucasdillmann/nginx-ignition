package api

import (
	"github.com/lucasdillmann/nginx-ignition/internal/api/core"
	"github.com/lucasdillmann/nginx-ignition/internal/api/domain"
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/container"
)

func Install() error {
	return container.Run(
		core.Install,
		domain.Install,
	)
}
