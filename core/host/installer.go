package host

import (
	"dillmann.com.br/nginx-ignition/core/certificate"
	"go.uber.org/dig"
)

func Install(container *dig.Container) error {
	return container.Provide(buildCommands)
}

func buildCommands(
	hostRepository Repository,
	certificateRepository certificate.Repository,
) (
	SaveCommand,
	DeleteCommand,
	ListCommand,
	GetCommand,
	GetAllEnabledCommand,
	ExistsCommand,
	ValidateBindingCommand,
) {
	serviceInstance := newService(&hostRepository, &certificateRepository)
	return serviceInstance.save,
		serviceInstance.deleteByID,
		serviceInstance.list,
		serviceInstance.getByID,
		serviceInstance.getAllEnabled,
		serviceInstance.existsByID,
		serviceInstance.validateBinding
}
