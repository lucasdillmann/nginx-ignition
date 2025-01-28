package host

import (
	"go.uber.org/dig"
)

func Install(container *dig.Container) error {
	return container.Provide(buildCommands)
}

func buildCommands(
	hostRepository Repository,
) (
	SaveCommand,
	DeleteCommand,
	ListCommand,
	GetCommand,
	GetAllEnabledCommand,
	ExistsCommand,
	ValidateBindingCommand,
) {
	serviceInstance := newService(&hostRepository)
	return serviceInstance.save,
		serviceInstance.deleteByID,
		serviceInstance.list,
		serviceInstance.getByID,
		serviceInstance.getAllEnabled,
		serviceInstance.existsByID,
		serviceInstance.validateBinding
}
