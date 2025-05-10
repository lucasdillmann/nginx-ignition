package core

import (
	"dillmann.com.br/nginx-ignition/core/access_list"
	"dillmann.com.br/nginx-ignition/core/certificate"
	"dillmann.com.br/nginx-ignition/core/common/broadcast"
	"dillmann.com.br/nginx-ignition/core/common/scheduler"
	"dillmann.com.br/nginx-ignition/core/host"
	"dillmann.com.br/nginx-ignition/core/integration"
	"dillmann.com.br/nginx-ignition/core/nginx"
	"dillmann.com.br/nginx-ignition/core/settings"
	"dillmann.com.br/nginx-ignition/core/stream"
	"dillmann.com.br/nginx-ignition/core/user"
	"go.uber.org/dig"
)

func Install(container *dig.Container) error {
	installers := []func(*dig.Container) error{
		broadcast.Install,
		scheduler.Install,
		settings.Install,
		user.Install,
		access_list.Install,
		certificate.Install,
		host.Install,
		integration.Install,
		stream.Install,
		nginx.Install,
	}

	for _, installer := range installers {
		if err := installer(container); err != nil {
			return err
		}
	}

	return nil
}
