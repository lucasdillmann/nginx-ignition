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
	driverResolver := func() ([]Driver, error) {
		var output []Driver
		if err := container.Invoke(func(drivers []Driver) {
			output = drivers
		}); err != nil {
			return nil, err
		}

		return output, nil
	}

	serviceInstance := newService(repository, driverResolver)

	return &Commands{
		Get:                 serviceInstance.getById,
		Save:                serviceInstance.save,
		Delete:              serviceInstance.deleteById,
		Exists:              serviceInstance.existsById,
		List:                serviceInstance.list,
		GetAvailableDrivers: serviceInstance.getAvailableDrivers,
		GetOption:           serviceInstance.getOptionById,
		GetOptionUrl:        serviceInstance.getOptionUrl,
		ListOptions:         serviceInstance.listOptions,
	}
}
