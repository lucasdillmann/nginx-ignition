package nginx

import (
	"github.com/lucasdillmann/nginx-ignition/internal/core/certificate"
	"github.com/lucasdillmann/nginx-ignition/internal/core/common/configuration"
	"github.com/lucasdillmann/nginx-ignition/internal/core/common/container"
	"github.com/lucasdillmann/nginx-ignition/internal/core/host"
	"github.com/lucasdillmann/nginx-ignition/internal/core/nginx/cfgfiles"
	"github.com/lucasdillmann/nginx-ignition/internal/core/settings"
	"github.com/lucasdillmann/nginx-ignition/internal/core/vpn"
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
