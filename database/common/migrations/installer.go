package migrations

import (
	"dillmann.com.br/nginx-ignition/core/common/container"
)

func Install() error {
	if err := container.Provide(newMigrations); err != nil {
		return err
	}

	return container.Run(registerStartup)
}
