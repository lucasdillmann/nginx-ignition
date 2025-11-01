package nginx

import (
	"dillmann.com.br/nginx-ignition/core/common/configuration"
	"dillmann.com.br/nginx-ignition/core/common/container"
	"dillmann.com.br/nginx-ignition/core/host"
	"dillmann.com.br/nginx-ignition/core/nginx/cfgfiles"
	"dillmann.com.br/nginx-ignition/core/settings"
)

func Install() error {
	if err := container.Run(cfgfiles.Install); err != nil {
		return err
	}

	if err := container.Provide(buildCommands); err != nil {
		return err
	}

	return container.Run(registerStartup, registerScheduledTask, registerShutdown)
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
		GetHostLogs:    serviceInstance.getHostLogs,
		GetMainLogs:    serviceInstance.getMainLogs,
		GetStatus:      serviceInstance.isRunning,
		GetConfigFiles: serviceInstance.getConfigFilesZipFile,
		GetMetadata:    serviceInstance.getMetadata,
		Reload:         serviceInstance.reload,
		Stop:           serviceInstance.stop,
		Start:          serviceInstance.start,
	}

	return serviceInstance, commands, nil
}
