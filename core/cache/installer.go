package cache

import (
	"dillmann.com.br/nginx-ignition/core/common/container"
)

func Install() error {
	return container.Provide(buildCommands)
}

func buildCommands(repository Repository) *Commands {
	serviceInstance := newService(repository)
	return &Commands{
		Save:        serviceInstance.save,
		Delete:      serviceInstance.deleteById,
		List:        serviceInstance.list,
		Get:         serviceInstance.findById,
		Exists:      serviceInstance.existsByID,
		GetAllInUse: serviceInstance.findAllInUse,
	}
}
