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
		Delete: serviceInstance.deleteById,
		Get:    serviceInstance.findById,
		Exists: serviceInstance.existsByID,
		List:   serviceInstance.list,
		Save:   serviceInstance.save,
	}
}
