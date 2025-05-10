package database

import (
	"dillmann.com.br/nginx-ignition/database/access_list"
	"dillmann.com.br/nginx-ignition/database/certificate"
	"dillmann.com.br/nginx-ignition/database/common/database"
	"dillmann.com.br/nginx-ignition/database/common/migrations"
	"dillmann.com.br/nginx-ignition/database/host"
	"dillmann.com.br/nginx-ignition/database/integration"
	"dillmann.com.br/nginx-ignition/database/settings"
	"dillmann.com.br/nginx-ignition/database/stream"
	"dillmann.com.br/nginx-ignition/database/user"
	"go.uber.org/dig"
)

func Install(container *dig.Container) error {
	installers := []func(*dig.Container) error{
		database.Install,
		migrations.Install,
	}

	for _, installer := range installers {
		if err := installer(container); err != nil {
			return err
		}
	}

	providers := []interface{}{
		access_list.New,
		host.New,
		user.New,
		settings.New,
		certificate.New,
		integration.New,
		stream.New,
	}

	for _, provider := range providers {
		if err := container.Provide(provider); err != nil {
			return err
		}
	}

	return nil
}
