package integration

import (
	"dillmann.com.br/nginx-ignition/core/common/container"
)

func Install() error {
	return container.Provide(buildCommands)
}

func buildCommands(
	repository Repository,
) *Commands {
	drivers := func() []Driver {
		return container.Get[[]Driver]()
	}

	serviceInstance := newService(repository, drivers)

	return &Commands{
		Get:                 serviceInstance.getByID,
		Save:                serviceInstance.save,
		Delete:              serviceInstance.deleteByID,
		Exists:              serviceInstance.existsByID,
		List:                serviceInstance.list,
		GetAvailableDrivers: serviceInstance.getAvailableDrivers,
		GetOption:           serviceInstance.getOptionByID,
		GetOptionURL:        serviceInstance.getOptionURL,
		ListOptions:         serviceInstance.listOptions,
	}
}
