package core

import (
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/broadcast"
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/container"
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/scheduler"
)

func Install() error {
	return container.Run(
		broadcast.Install,
		scheduler.Install,
	)
}
