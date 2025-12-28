package settings

import (
	"dillmann.com.br/nginx-ignition/core/binding"
	"dillmann.com.br/nginx-ignition/core/common/container"
	"dillmann.com.br/nginx-ignition/core/common/scheduler"
)

func Install() error {
	return container.Provide(buildCommands)
}

func buildCommands(
	repository Repository,
	bindingCommands *binding.Commands,
	sched *scheduler.Scheduler,
) *Commands {
	serviceInstance := newService(repository, bindingCommands, sched)
	return &Commands{
		Get:  serviceInstance.get,
		Save: serviceInstance.save,
	}
}
