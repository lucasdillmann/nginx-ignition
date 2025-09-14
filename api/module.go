package api

import (
	"go.uber.org/dig"

	"dillmann.com.br/nginx-ignition/api/access_list"
	"dillmann.com.br/nginx-ignition/api/certificate"
	"dillmann.com.br/nginx-ignition/api/common/server"
	"dillmann.com.br/nginx-ignition/api/frontend"
	"dillmann.com.br/nginx-ignition/api/host"
	"dillmann.com.br/nginx-ignition/api/integration"
	"dillmann.com.br/nginx-ignition/api/nginx"
	"dillmann.com.br/nginx-ignition/api/settings"
	"dillmann.com.br/nginx-ignition/api/stream"
	"dillmann.com.br/nginx-ignition/api/user"
)

func Install(container *dig.Container) error {
	if err := server.Install(container); err != nil {
		return err
	}

	installFunctions := []interface{}{
		settings.Install,
		access_list.Install,
		certificate.Install,
		user.Install,
		host.Install,
		integration.Install,
		nginx.Install,
		stream.Install,
		frontend.Install,
	}

	for _, installFunc := range installFunctions {
		if err := container.Invoke(installFunc); err != nil {
			return err
		}
	}

	return nil
}
