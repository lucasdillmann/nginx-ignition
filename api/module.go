package api

import (
	"dillmann.com.br/nginx-ignition/api/common/server"
	"go.uber.org/dig"
)

func Install(container *dig.Container) error {
	if err := server.Install(container); err != nil {
		return err
	}

	return nil
}
