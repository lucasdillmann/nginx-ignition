package core

import (
	"dillmann.com.br/nginx-ignition/core/access_list"
	"dillmann.com.br/nginx-ignition/core/certificate"
	"dillmann.com.br/nginx-ignition/core/common/broadcast"
	"dillmann.com.br/nginx-ignition/core/common/scheduler"
	"dillmann.com.br/nginx-ignition/core/host"
	"dillmann.com.br/nginx-ignition/core/nginx"
	"dillmann.com.br/nginx-ignition/core/settings"
	"dillmann.com.br/nginx-ignition/core/user"
	"go.uber.org/dig"
)

func Install(container *dig.Container) error {
	if err := broadcast.Install(container); err != nil {
		return err
	}

	if err := scheduler.Install(container); err != nil {
		return err
	}

	if err := settings.Install(container); err != nil {
		return err
	}

	if err := user.Install(container); err != nil {
		return err
	}

	if err := access_list.Install(container); err != nil {
		return err
	}

	if err := certificate.Install(container); err != nil {
		return err
	}

	if err := host.Install(container); err != nil {
		return err
	}

	if err := nginx.Install(container); err != nil {
		return err
	}

	return nil
}
