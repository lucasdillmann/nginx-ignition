package vpn

import (
	"dillmann.com.br/nginx-ignition/core/common/configuration"
	"dillmann.com.br/nginx-ignition/core/common/container"
)

func Install() error {
	return container.Provide(newCommands)
}

func newCommands(cfg *configuration.Configuration, repository Repository) Commands {
	return newService(cfg, repository, func() []Driver {
		return container.Get[[]Driver]()
	})
}
