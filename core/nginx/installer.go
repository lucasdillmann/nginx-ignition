package nginx

import (
	"go.uber.org/dig"

	"dillmann.com.br/nginx-ignition/core/common/configuration"
	"dillmann.com.br/nginx-ignition/core/host"
	"dillmann.com.br/nginx-ignition/core/nginx/cfgfiles"
	"dillmann.com.br/nginx-ignition/core/settings"
)

func Install(container *dig.Container) error {
	if err := cfgfiles.Install(container); err != nil {
		return err
	}

	if err := container.Provide(buildCommands); err != nil {
		return err
	}

	if err := container.Invoke(registerStartup); err != nil {
		return err
	}

	if err := container.Invoke(registerScheduledTask); err != nil {
		return err
	}

	return container.Invoke(registerShutdown)
}

func buildCommands(
	configuration *configuration.Configuration,
	hostRepository host.Repository,
	settingsRepository settings.Repository,
	configFilesManager *cfgfiles.Facade,
) (*service, *Commands, error) {
	serviceInstance, err := newService(configuration, settingsRepository, hostRepository, configFilesManager)
	if err != nil {
		return nil, nil, err
	}

	commands := &Commands{
		GetHostLogs: serviceInstance.getHostLogs,
		GetMainLogs: serviceInstance.getMainLogs,
		GetStatus:   serviceInstance.isRunning,
		Reload:      serviceInstance.reload,
		Stop:        serviceInstance.stop,
		Start:       serviceInstance.start,
	}

	return serviceInstance, commands, nil
}
