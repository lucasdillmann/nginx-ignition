package domain

import (
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/container"
	"github.com/lucasdillmann/nginx-ignition/internal/vpn/domain/netbird"
	"github.com/lucasdillmann/nginx-ignition/internal/vpn/domain/tailscale"
)

func Install() error {
	return container.Run(
		netbird.Install,
		tailscale.Install,
	)
}
