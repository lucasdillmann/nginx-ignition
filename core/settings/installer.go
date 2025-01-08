package settings

import "go.uber.org/dig"

func Install(container *dig.Container) error {
	return container.Provide(buildCommands)
}

func buildCommands(repository *Repository) (GetCommand, SaveCommand) {
	serviceInstance := &service{repository}
	return serviceInstance.Get, serviceInstance.Save
}
