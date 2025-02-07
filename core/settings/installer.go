package settings

import (
	"dillmann.com.br/nginx-ignition/core/common/scheduler"
	"dillmann.com.br/nginx-ignition/core/host"
	"go.uber.org/dig"
)

func Install(container *dig.Container) error {
	return container.Provide(buildCommands)
}

func buildCommands(
	repository Repository,
	validateBindingsCommand host.ValidateBindingCommand,
	scheduler *scheduler.Scheduler,
) (GetCommand, SaveCommand) {
	serviceInstance := newService(&repository, &validateBindingsCommand, scheduler)
	return serviceInstance.get, serviceInstance.save
}
