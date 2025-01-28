package access_list

import (
	"dillmann.com.br/nginx-ignition/core/host"
	"go.uber.org/dig"
)

func Install(container *dig.Container) error {
	return container.Provide(buildCommands)
}

func buildCommands(
	accessListRepository Repository,
	hostRepository host.Repository,
) (
	GetCommand,
	DeleteCommand,
	ListCommand,
	SaveCommand,
) {
	serviceInstance := newService(&accessListRepository, &hostRepository)
	return serviceInstance.findById, serviceInstance.deleteById, serviceInstance.list, serviceInstance.save
}
