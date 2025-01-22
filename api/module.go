package api

import (
	"dillmann.com.br/nginx-ignition/api/access_list"
	"dillmann.com.br/nginx-ignition/api/certificate"
	"dillmann.com.br/nginx-ignition/api/common/server"
	"dillmann.com.br/nginx-ignition/api/frontend"
	"dillmann.com.br/nginx-ignition/api/host"
	"dillmann.com.br/nginx-ignition/api/nginx"
	"dillmann.com.br/nginx-ignition/api/settings"
	"dillmann.com.br/nginx-ignition/api/user"
	"go.uber.org/dig"
)

func Install(container *dig.Container) error {
	if err := server.Install(container); err != nil {
		return err
	}

	if err := container.Invoke(settings.Install); err != nil {
		return err
	}

	if err := container.Invoke(access_list.Install); err != nil {
		return err
	}

	if err := container.Invoke(certificate.Install); err != nil {
		return err
	}

	if err := container.Invoke(user.Install); err != nil {
		return err
	}

	if err := container.Invoke(host.Install); err != nil {
		return err
	}

	if err := container.Invoke(nginx.Install); err != nil {
		return err
	}

	if err := container.Invoke(frontend.Install); err != nil {
		return err
	}

	return nil
}
