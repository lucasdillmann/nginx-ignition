package user

import (
	"dillmann.com.br/nginx-ignition/core/common/configuration"
	"dillmann.com.br/nginx-ignition/core/common/container"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
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
	i18nCommands i18n.Commands,
) (*service, Commands) {
	serviceInstance := newService(repository, cfg, i18nCommands)
	return serviceInstance, serviceInstance
}
