package core

import (
	"github.com/lucasdillmann/nginx-ignition/internal/api/core/server"
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/container"
)

func Install() error {
	return container.Run(
		server.Install,
	)
}
