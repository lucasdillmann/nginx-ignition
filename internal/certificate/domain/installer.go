package domain

import (
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/container"
	"github.com/lucasdillmann/nginx-ignition/internal/certificate/domain/custom"
	"github.com/lucasdillmann/nginx-ignition/internal/certificate/domain/external"
	"github.com/lucasdillmann/nginx-ignition/internal/certificate/domain/letsencrypt"
	"github.com/lucasdillmann/nginx-ignition/internal/certificate/domain/selfsigned"
)

func Install() error {
	return container.Run(
		custom.Install,
		external.Install,
		letsencrypt.Install,
		selfsigned.Install,
	)
}
