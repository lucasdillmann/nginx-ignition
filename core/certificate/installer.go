package certificate

import (
	"dillmann.com.br/nginx-ignition/core/host"
	"dillmann.com.br/nginx-ignition/core/settings"
	"go.uber.org/dig"
)

func Install(container *dig.Container) error {
	if err := container.Provide(buildCommands); err != nil {
		return err
	}

	return container.Invoke(registerScheduledTask)
}

func buildCommands(
	container *dig.Container,
	hostRepository host.Repository,
	certificateRepository Repository,
	settingsRepository settings.Repository,
) (*Commands, *service) {
	providerResolver := func() ([]Provider, error) {
		var output []Provider
		if err := container.Invoke(func(providers []Provider) {
			output = providers
		}); err != nil {
			return nil, err
		}

		return output, nil
	}

	serviceInstance := newService(certificateRepository, hostRepository, settingsRepository, providerResolver)
	commands := &Commands{
		AvailableProviders: serviceInstance.availableProviders,
		Delete:             serviceInstance.deleteById,
		Get:                serviceInstance.getById,
		List:               serviceInstance.list,
		Renew:              serviceInstance.renew,
	}

	return commands, serviceInstance
}
