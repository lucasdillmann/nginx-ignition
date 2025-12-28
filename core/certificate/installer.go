package certificate

import (
	"dillmann.com.br/nginx-ignition/core/common/container"
)

func Install() error {
	if err := container.Provide(buildCommands); err != nil {
		return err
	}

	return container.Run(registerScheduledTask)
}

func buildCommands(repository Repository) (*Commands, *service) {
	providers := func() []Provider {
		return container.Get[[]Provider]()
	}

	serviceInstance := newService(repository, providers)
	commands := &Commands{
		AvailableProviders: serviceInstance.availableProviders,
		Delete:             serviceInstance.deleteByID,
		Get:                serviceInstance.getByID,
		Exists:             serviceInstance.existsByID,
		List:               serviceInstance.list,
		Renew:              serviceInstance.renew,
		Issue:              serviceInstance.issue,
	}

	return commands, serviceInstance
}
