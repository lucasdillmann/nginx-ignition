package migrations

import (
	"github.com/lucasdillmann/nginx-ignition/internal/core/common/container"
)

func Install() error {
	if err := container.Provide(New); err != nil {
		return err
	}

	return container.Run(registerStartup)
}
