package server

import (
	"dillmann.com.br/nginx-ignition/core/common/container"
	"dillmann.com.br/nginx-ignition/core/host"
	"dillmann.com.br/nginx-ignition/core/settings"
)

func Install() error {
	if err := container.Provide(buildCommands); err != nil {
		return err
	}

	return container.Run(registerScheduledTask)
}

func buildCommands(
	hostRepository host.Repository,
	certificateRepository Repository,
	settingsRepository settings.Repository,
) (*Commands, *service) {
	providers := func() []Provider {
		return container.Get[[]Provider]()
	}

	serviceInstance := newService(certificateRepository, hostRepository, settingsRepository, providers)
	commands := &Commands{
		AvailableProviders: serviceInstance.availableProviders,
		Delete:             serviceInstance.deleteById,
		Get:                serviceInstance.getById,
		List:               serviceInstance.list,
		Renew:              serviceInstance.renew,
		Issue:              serviceInstance.issue,
	}

	return commands, serviceInstance
}
