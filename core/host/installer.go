package host

import (
	"dillmann.com.br/nginx-ignition/core/common/container"
	"dillmann.com.br/nginx-ignition/core/integration"
)

func Install() error {
	return container.Provide(buildCommands)
}

func buildCommands(
	hostRepository Repository,
	integrationCommands *integration.Commands,
) *Commands {
	serviceInstance := newService(hostRepository, integrationCommands)
	return &Commands{
		Save:            serviceInstance.save,
		Delete:          serviceInstance.deleteByID,
		List:            serviceInstance.list,
		Get:             serviceInstance.getByID,
		GetAllEnabled:   serviceInstance.getAllEnabled,
		Exists:          serviceInstance.existsByID,
		ValidateBinding: serviceInstance.validateBinding,
	}
}
