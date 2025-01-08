package database

import (
	"dillmann.com.br/nginx-ignition/database/access_list_repository"
	"dillmann.com.br/nginx-ignition/database/common/database"
	"dillmann.com.br/nginx-ignition/database/host_repository"
	"dillmann.com.br/nginx-ignition/database/user_repository"
	"go.uber.org/dig"
)

func Install(container *dig.Container) error {
	if err := database.Install(container); err != nil {
		return err
	}

	if err := container.Provide(access_list_repository.New); err != nil {
		return err
	}

	if err := container.Provide(host_repository.New); err != nil {
		return err
	}

	if err := container.Provide(user_repository.New); err != nil {
		return err
	}

	return nil
}
