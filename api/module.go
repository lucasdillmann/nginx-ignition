package api

import (
	"dillmann.com.br/nginx-ignition/api/common/server"
	"dillmann.com.br/nginx-ignition/api/settings_api"
	"go.uber.org/dig"
)

func Install(container *dig.Container) error {
	if err := server.Install(container); err != nil {
		return err
	}

	if err := container.Invoke(settings_api.Install); err != nil {
		return err
	}

	return nil
}
