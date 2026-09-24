package vpn

import (
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/container"
	"github.com/lucasdillmann/nginx-ignition/internal/vpn/domain"
)

func Install() error {
	return container.Run(
		domain.Install,
	)
}
