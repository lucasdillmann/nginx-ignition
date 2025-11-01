package settings

import (
	"dillmann.com.br/nginx-ignition/core/common/container"
	"dillmann.com.br/nginx-ignition/core/common/scheduler"
	"dillmann.com.br/nginx-ignition/core/host"
)

func Install() error {
	return container.Provide(buildCommands)
}

func buildCommands(
	repository Repository,
	hostCommands *host.Commands,
	scheduler *scheduler.Scheduler,
) *Commands {
	serviceInstance := newService(repository, hostCommands, scheduler)
	return &Commands{
		Get:  serviceInstance.get,
		Save: serviceInstance.save,
	}
}
