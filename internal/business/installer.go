package business

import (
	"github.com/lucasdillmann/nginx-ignition/internal/business/core"
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/container"
	"github.com/lucasdillmann/nginx-ignition/internal/business/domain"
)

func Install() error {
	return container.Run(
		core.Install,
		domain.Install,
	)
}
