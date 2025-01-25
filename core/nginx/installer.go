package nginx

import (
	"dillmann.com.br/nginx-ignition/core/common/configuration"
	"dillmann.com.br/nginx-ignition/core/host"
	"dillmann.com.br/nginx-ignition/core/nginx/cfgfiles"
	"dillmann.com.br/nginx-ignition/core/settings"
	"go.uber.org/dig"
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

	return container.Invoke(registerShutdown)
}

func buildCommands(
	configuration *configuration.Configuration,
	hostRepository host.Repository,
	settingsRepository settings.Repository,
	configFilesManager *cfgfiles.Facade,
) (
	*service,
	GetHostLogsCommand,
	GetMainLogsCommand,
	GetStatusCommand,
	ReloadCommand,
	StopCommand,
	StartCommand,
	error,
) {
	serviceInstance, err := newService(configuration, &settingsRepository, &hostRepository, configFilesManager)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, err
	}

	return serviceInstance,
		serviceInstance.getHostLogs,
		serviceInstance.getMainLogs,
		serviceInstance.isRunning,
		serviceInstance.reload,
		serviceInstance.stop,
		serviceInstance.start,
		nil
}
