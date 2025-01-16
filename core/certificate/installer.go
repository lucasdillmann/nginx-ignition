package certificate

import (
	"dillmann.com.br/nginx-ignition/core/host"
	"dillmann.com.br/nginx-ignition/core/settings"
	"go.uber.org/dig"
)

func Install(container *dig.Container) error {
	return container.Provide(buildCommands)
}

func buildCommands(
	container *dig.Container,
	hostRepository host.Repository,
	certificateRepository Repository,
	settingsRepository settings.Repository,
) (
	IssueCommand,
	RenewCommand,
	DeleteCommand,
	GetCommand,
	ListCommand,
	AvailableProvidersCommand,
) {
	serviceInstance := newService(container, &certificateRepository, &hostRepository, &settingsRepository)
	return serviceInstance.issue,
		serviceInstance.renew,
		serviceInstance.deleteById,
		serviceInstance.getById,
		serviceInstance.list,
		serviceInstance.availableProviders
}
