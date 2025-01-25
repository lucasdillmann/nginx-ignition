package certificate

import (
	"dillmann.com.br/nginx-ignition/core/common/scheduler"
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

func registerScheduledTask(
	service *service,
	settingsRepository settings.Repository,
	scheduler *scheduler.Scheduler,
) error {
	task := autoRenewTask{service, &settingsRepository}
	return scheduler.Register(&task)
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
	serviceInstance := newService(container, &certificateRepository, &hostRepository, &settingsRepository)
	return serviceInstance,
		serviceInstance.issue,
		serviceInstance.renew,
		serviceInstance.deleteById,
		serviceInstance.getById,
		serviceInstance.list,
		serviceInstance.availableProviders
}
