package vpn

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
		Get:                 serviceInstance.getById,
		Save:                serviceInstance.save,
		Delete:              serviceInstance.deleteById,
		Exists:              serviceInstance.existsById,
		List:                serviceInstance.list,
		GetAvailableDrivers: serviceInstance.getAvailableDrivers,
	}
}
