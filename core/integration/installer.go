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
) *Commands {
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

	return &Commands{
		ConfigureById:    serviceInstance.configureById,
		GetById:          serviceInstance.getById,
		GetOptionById:    serviceInstance.getOptionById,
		GetOptionUrlById: serviceInstance.getOptionUrl,
		ListOptions:      serviceInstance.listOptions,
		List:             serviceInstance.list,
	}
}
