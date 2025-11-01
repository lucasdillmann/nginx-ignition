package user

import (
	"dillmann.com.br/nginx-ignition/core/common/configuration"
	"dillmann.com.br/nginx-ignition/core/common/container"
)

func Install() error {
	err := container.Provide(buildCommands)
	if err != nil {
		return err
	}

	return container.Run(registerStartup)
}

func buildCommands(
	repository Repository,
	configuration *configuration.Configuration,
) (
	*service,
	*Commands,
) {
	serviceInstance := &service{repository, configuration}
	commands := &Commands{
		Authenticate:        serviceInstance.authenticate,
		Delete:              serviceInstance.deleteById,
		Get:                 serviceInstance.getById,
		GetCount:            serviceInstance.count,
		GetStatus:           serviceInstance.isEnabled,
		List:                serviceInstance.list,
		Save:                serviceInstance.save,
		UpdatePassword:      serviceInstance.changePassword,
		OnboardingCompleted: serviceInstance.isOnboardingCompleted,
	}

	return serviceInstance, commands
}
