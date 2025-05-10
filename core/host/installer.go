package host

import (
	"go.uber.org/dig"
)

func Install(container *dig.Container) error {
	return container.Provide(buildCommands)
}

func buildCommands(
	hostRepository Repository,
) *Commands {
	serviceInstance := newService(hostRepository)
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
