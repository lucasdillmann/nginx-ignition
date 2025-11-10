package backup

import (
	"dillmann.com.br/nginx-ignition/core/common/container"
)

func Install() error {
	return container.Provide(buildCommands)
}

func buildCommands(repository Repository) *Commands {
	serviceInstance := newService(repository)
	return &Commands{
		Get: serviceInstance.get,
	}
}
