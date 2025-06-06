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
) *Commands {
	serviceInstance := newService(accessListRepository, hostRepository)
	return &Commands{
		Delete: serviceInstance.deleteById,
		Get:    serviceInstance.findById,
		List:   serviceInstance.list,
		Save:   serviceInstance.save,
	}
}
