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
	DeleteByIdCommand,
	GetByIdCommand,
	GetCountCommand,
	GetStatusCommand,
	ListCommand,
	SaveCommand,
	ChangePasswordCommand,
) {
	serviceInstance := &service{&repository, configuration}
	return serviceInstance.authenticate,
		serviceInstance.deleteById,
		serviceInstance.getById,
		serviceInstance.count,
		serviceInstance.isEnabled,
		serviceInstance.list,
		serviceInstance.save,
		serviceInstance.changePassword
}
