package database

import (
	"dillmann.com.br/nginx-ignition/core/common/container"
)

func Install() error {
	if err := container.Provide(newDatabase); err != nil {
		return err
	}

	return container.Run(registerStartup, registerShutdown, registerHealthCheck)
}
