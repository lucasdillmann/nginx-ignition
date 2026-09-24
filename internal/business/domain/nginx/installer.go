package nginx

import (
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/configuration"
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/container"
	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/certificate"
	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/host"
	cfgfiles2 "github.com/lucasdillmann/nginx-ignition/internal/business/domain/nginx/cfgfiles"
	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/settings"
	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/vpn"
)

func Install() error {
	if err := container.Run(cfgfiles2.Install); err != nil {
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
	configFilesManager *cfgfiles2.Facade,
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
