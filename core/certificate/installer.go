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
) (
	*service,
	IssueCommand,
	RenewCommand,
	DeleteCommand,
	GetCommand,
	ListCommand,
	AvailableProvidersCommand,
) {
	providerResolver := func() ([]Provider, error) {
		var output []Provider
		if err := container.Invoke(func(providers []Provider) {
			output = providers
		}); err != nil {
			return nil, err
		}

		return output, nil
	}

	serviceInstance := newService(&certificateRepository, &hostRepository, &settingsRepository, providerResolver)

	return serviceInstance,
		serviceInstance.issue,
		serviceInstance.renew,
		serviceInstance.deleteById,
		serviceInstance.getById,
		serviceInstance.list,
		serviceInstance.availableProviders
}
