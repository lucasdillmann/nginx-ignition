package api

import (
	"dillmann.com.br/nginx-ignition/api/access_list"
	"dillmann.com.br/nginx-ignition/api/common/server"
	"dillmann.com.br/nginx-ignition/api/host"
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

	if err := container.Invoke(user.Install); err != nil {
		return err
	}

	if err := container.Invoke(host.Install); err != nil {
		return err
	}

	return nil
}
