package host

import (
	"dillmann.com.br/nginx-ignition/core/accesslist"
	"dillmann.com.br/nginx-ignition/core/cache"
	"dillmann.com.br/nginx-ignition/core/common/container"
	"dillmann.com.br/nginx-ignition/core/integration"
	"dillmann.com.br/nginx-ignition/core/vpn"
)

func Install() error {
	return container.Provide(buildCommands)
}

func buildCommands(
	hostRepository Repository,
	integrationCommands *integration.Commands,
	vpnCommands *vpn.Commands,
	accessListCommands *accesslist.Commands,
	cacheCommands *cache.Commands,
) *Commands {
	serviceInstance := newService(
		hostRepository,
		integrationCommands,
		vpnCommands,
		accessListCommands,
		cacheCommands,
	)
	return &Commands{
		Save:            serviceInstance.save,
		Delete:          serviceInstance.deleteByID,
		List:            serviceInstance.list,
		Get:             serviceInstance.getByID,
		GetAllEnabled:   serviceInstance.getAllEnabled,
		Exists:          serviceInstance.existsByID,
		ValidateBinding: serviceInstance.validateBinding,
	}
}
