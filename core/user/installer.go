package user

import (
	"dillmann.com.br/nginx-ignition/core/common/configuration"
	"dillmann.com.br/nginx-ignition/core/common/container"
)

func Install() error {
	err := container.Provide(newCommands)
	if err != nil {
		return err
	}

	return container.Run(registerStartup)
}

func newCommands(
	repository Repository,
	cfg *configuration.Configuration,
) (*service, Commands) {
	serviceInstance := newService(repository, cfg)
	return serviceInstance, serviceInstance
}
