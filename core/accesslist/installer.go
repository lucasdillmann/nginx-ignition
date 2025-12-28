package accesslist

import (
	"dillmann.com.br/nginx-ignition/core/common/container"
)

func Install() error {
	return container.Provide(buildCommands)
}

func buildCommands(repository Repository) *Commands {
	serviceInstance := newService(repository)
	return &Commands{
		Delete: serviceInstance.deleteByID,
		Get:    serviceInstance.findByID,
		GetAll: serviceInstance.findAll,
		List:   serviceInstance.list,
		Save:   serviceInstance.save,
		Exists: serviceInstance.existsByID,
	}
}
