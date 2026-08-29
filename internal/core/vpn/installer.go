package vpn

import (
	"nginx-ignition/internal/core/common/configuration"
	"nginx-ignition/internal/core/common/container"
)

func Install() error {
	return container.Provide(newCommands)
}

func newCommands(cfg *configuration.Configuration, repository Repository) Commands {
	return newService(cfg, repository, func() []Driver {
		return container.Get[[]Driver]()
	})
}
