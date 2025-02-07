package integration

import (
	"go.uber.org/dig"
)

func Install(container *dig.Container) error {
	return container.Provide(buildCommands)
}

func buildCommands(
	container *dig.Container,
	repository Repository,
) (
	ConfigureByIdCommand,
	GetByIdCommand,
	GetOptionByIdCommand,
	GetOptionUrlByIdCommand,
	ListOptionsCommand,
	ListCommand,
) {
	adapterResolver := func() ([]Adapter, error) {
		var output []Adapter
		if err := container.Invoke(func(adapters []Adapter) {
			output = adapters
		}); err != nil {
			return nil, err
		}

		return output, nil
	}

	serviceInstance := newService(repository, adapterResolver)

	return serviceInstance.configureById,
		serviceInstance.getById,
		serviceInstance.getOptionById,
		serviceInstance.getOptionUrl,
		serviceInstance.listOptions,
		serviceInstance.list
}
