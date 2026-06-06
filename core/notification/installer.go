package notification

import (
	"dillmann.com.br/nginx-ignition/core/common/container"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
	"dillmann.com.br/nginx-ignition/core/user"
)

func Install() error {
	if err := container.Provide(newCommands); err != nil {
		return err
	}

	return container.Run(registerScheduledTask)
}

func newCommands(
	repository Repository,
	userCommands user.Commands,
	i18nCommands i18n.Commands,
) (Commands, *service) {
	providers := func() []Provider {
		return container.Get[[]Provider]()
	}

	serviceInstance := newService(
		repository,
		userCommands,
		i18nCommands,
		providers,
	)

	return serviceInstance, serviceInstance
}
