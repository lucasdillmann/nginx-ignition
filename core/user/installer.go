package user

import (
	"dillmann.com.br/nginx-ignition/core/common/configuration"
	"go.uber.org/dig"
)

func Install(container *dig.Container) error {
	return container.Provide(buildCommands)
}

func buildCommands(
	repository Repository,
	configuration *configuration.Configuration,
) (
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
	return serviceInstance.authenticate,
		serviceInstance.deleteById,
		serviceInstance.getById,
		serviceInstance.count,
		serviceInstance.isEnabled,
		serviceInstance.list,
		serviceInstance.save,
		serviceInstance.changePassword,
		serviceInstance.isOnboardingCompleted
}
