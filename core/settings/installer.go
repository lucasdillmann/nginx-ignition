package settings

import (
	"dillmann.com.br/nginx-ignition/core/host"
	"go.uber.org/dig"
)

func Install(container *dig.Container) error {
	return container.Provide(buildCommands)
}

func buildCommands(
	repository Repository,
	validateBindingsCommand host.ValidateBindingCommand,
) (GetCommand, SaveCommand) {
	serviceInstance := newService(&repository, &validateBindingsCommand)
	return serviceInstance.get, serviceInstance.save
}
