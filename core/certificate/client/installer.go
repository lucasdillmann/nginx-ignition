package client

import (
	"dillmann.com.br/nginx-ignition/core/common/container"
)

func Install() error {
	return container.Provide(buildCommands)
}

func buildCommands(repository Repository) (*Commands, *service) {
	serviceInstance := newService(repository)
	commands := &Commands{
		List:         serviceInstance.list,
		Create:       nil,
		Delete:       serviceInstance.deleteById,
		Get:          serviceInstance.getById,
		Update:       nil,
		ReplaceCA:    nil,
		UpdateClient: nil,
		CreateClient: nil,
		DeleteClient: nil,
	}

	return commands, serviceInstance
}
