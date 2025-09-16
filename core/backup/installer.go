package backup

import (
	"go.uber.org/dig"
)

func Install(container *dig.Container) error {
	return container.Provide(buildCommands)
}

func buildCommands(repository Repository) *Commands {
	serviceInstance := newService(repository)
	return &Commands{
		Get: serviceInstance.get,
	}
}
