package api

import (
	"dillmann.com.br/nginx-ignition/api/access_list_api"
	"dillmann.com.br/nginx-ignition/api/common/server"
	"dillmann.com.br/nginx-ignition/api/settings_api"
	"dillmann.com.br/nginx-ignition/api/user_api"
	"go.uber.org/dig"
)

func Install(container *dig.Container) error {
	if err := server.Install(container); err != nil {
		return err
	}

	if err := container.Invoke(settings_api.Install); err != nil {
		return err
	}

	if err := container.Invoke(access_list_api.Install); err != nil {
		return err
	}

	if err := container.Invoke(user_api.Install); err != nil {
		return err
	}

	return nil
}
