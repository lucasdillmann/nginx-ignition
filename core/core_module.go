package core

import (
	"dillmann.com.br/nginx-ignition/core/access_list"
	"dillmann.com.br/nginx-ignition/core/user"
	"go.uber.org/dig"
)

func InstallBeans(container *dig.Container) error {
	if err := user.InstallBeans(container); err != nil {
		return err
	}

	if err := access_list.InstallBeans(container); err != nil {
		return err
	}

	return nil
}
