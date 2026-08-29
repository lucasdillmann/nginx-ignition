package nginx

import (
	"nginx-ignition/internal/core/certificate"
	"nginx-ignition/internal/core/common/configuration"
	"nginx-ignition/internal/core/common/container"
	"nginx-ignition/internal/core/host"
	"nginx-ignition/internal/core/nginx/cfgfiles"
	"nginx-ignition/internal/core/settings"
	"nginx-ignition/internal/core/vpn"
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
	certificateCommands certificate.Commands,
) (*service, Commands, error) {
	serviceInstance, err := newService(
		cfg,
		hostCommands,
		configFilesManager,
		vpnCommands,
		settingsCommands,
		certificateCommands,
	)
	if err != nil {
		return nil, nil, err
	}

	return serviceInstance, serviceInstance, nil
}
