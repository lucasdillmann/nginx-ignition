package certificate

import (
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/container"
	"github.com/lucasdillmann/nginx-ignition/internal/certificate/domain"
)

func Install() error {
	return container.Run(
		domain.Install,
	)
}
