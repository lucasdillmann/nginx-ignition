package database

import (
	"dillmann.com.br/nginx-ignition/database/access_list"
	"dillmann.com.br/nginx-ignition/database/certificate"
	"dillmann.com.br/nginx-ignition/database/common/database"
	"dillmann.com.br/nginx-ignition/database/common/migrations"
	"dillmann.com.br/nginx-ignition/database/host"
	"dillmann.com.br/nginx-ignition/database/integration"
	"dillmann.com.br/nginx-ignition/database/settings"
	"dillmann.com.br/nginx-ignition/database/user"
	"go.uber.org/dig"
)

func Install(container *dig.Container) error {
	if err := database.Install(container); err != nil {
		return err
	}

	if err := migrations.Install(container); err != nil {
		return err
	}

	if err := container.Provide(access_list.New); err != nil {
		return err
	}

	if err := container.Provide(host.New); err != nil {
		return err
	}

	if err := container.Provide(user.New); err != nil {
		return err
	}

	if err := container.Provide(settings.New); err != nil {
		return err
	}

	if err := container.Provide(certificate.New); err != nil {
		return err
	}

	if err := container.Provide(integration.New); err != nil {
		return err
	}

	return nil
}
