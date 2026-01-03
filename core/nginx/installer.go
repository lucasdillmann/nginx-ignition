package nginx

import (
	"dillmann.com.br/nginx-ignition/core/common/configuration"
	"dillmann.com.br/nginx-ignition/core/common/container"
	"dillmann.com.br/nginx-ignition/core/host"
	"dillmann.com.br/nginx-ignition/core/nginx/cfgfiles"
	"dillmann.com.br/nginx-ignition/core/settings"
	"dillmann.com.br/nginx-ignition/core/vpn"
)

func Install() error {
	if err := container.Run(cfgfiles.Install); err != nil {
		return err
	}

	if err := container.Provide(newCommands); err != nil {
		return err
	}

	return container.Run(registerStartup, registerScheduledTask, registerShutdown)
}

func newCommands(
	cfg *configuration.Configuration,
	hostCommands host.Commands,
	configFilesManager *cfgfiles.Facade,
	vpnCommands vpn.Commands,
	settingsCommands settings.Commands,
) (*service, Commands, error) {
	serviceInstance, err := newService(
		cfg,
		hostCommands,
		configFilesManager,
		vpnCommands,
		settingsCommands,
	)
	if err != nil {
		return nil, nil, err
	}

	return serviceInstance, serviceInstance, nil
}
