package user

import (
	"dillmann.com.br/nginx-ignition/core/common/configuration"
	"go.uber.org/dig"
)

func Install(container *dig.Container) error {
	err := container.Provide(buildCommands)
	if err != nil {
		return err
	}

	return container.Invoke(registerStartup)
}

func buildCommands(
	repository Repository,
	configuration *configuration.Configuration,
) (
	*service,
	AuthenticateCommand,
	DeleteCommand,
	GetCommand,
	GetCountCommand,
	GetStatusCommand,
	ListCommand,
	SaveCommand,
	UpdatePasswordCommand,
	OnboardingCompletedCommand,
) {
	serviceInstance := &service{&repository, configuration}
	return serviceInstance,
		serviceInstance.authenticate,
		serviceInstance.deleteById,
		serviceInstance.getById,
		serviceInstance.count,
		serviceInstance.isEnabled,
		serviceInstance.list,
		serviceInstance.save,
		serviceInstance.changePassword,
		serviceInstance.isOnboardingCompleted
}
