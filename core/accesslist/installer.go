package accesslist

import (
	"dillmann.com.br/nginx-ignition/core/common/container"
	"dillmann.com.br/nginx-ignition/core/host"
)

func Install() error {
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
