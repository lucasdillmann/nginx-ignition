package vpn

import (
	"dillmann.com.br/nginx-ignition/core/common/configuration"
	"dillmann.com.br/nginx-ignition/core/common/container"
)

func Install() error {
	return container.Provide(buildCommands)
}

func buildCommands(cfg *configuration.Configuration, repository Repository) *Commands {
	serviceInstance := newService(cfg, repository, func() []Driver {
		return container.Get[[]Driver]()
	})

	return &Commands{
		Get:                 serviceInstance.getById,
		Save:                serviceInstance.save,
		Delete:              serviceInstance.deleteById,
		Exists:              serviceInstance.existsById,
		List:                serviceInstance.list,
		GetAvailableDrivers: serviceInstance.getAvailableDrivers,
		Start:               serviceInstance.start,
		Reload:              serviceInstance.reload,
		Stop:                serviceInstance.stop,
	}
}
