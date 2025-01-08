package core

import (
	"dillmann.com.br/nginx-ignition/core/access_list"
	"dillmann.com.br/nginx-ignition/core/settings"
	"dillmann.com.br/nginx-ignition/core/user"
	"go.uber.org/dig"
)

func Install(container *dig.Container) error {
	if err := user.Install(container); err != nil {
		return err
	}

	if err := access_list.Install(container); err != nil {
		return err
	}

	if err := settings.Install(container); err != nil {
		return err
	}

	return nil
}
