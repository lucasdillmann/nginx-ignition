package cfgfiles

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"dillmann.com.br/nginx-ignition/core/access_list"
	"dillmann.com.br/nginx-ignition/core/certificate"
	"dillmann.com.br/nginx-ignition/core/common/configuration"
	"dillmann.com.br/nginx-ignition/core/common/log"
	"dillmann.com.br/nginx-ignition/core/host"
	"dillmann.com.br/nginx-ignition/core/integration"
	"dillmann.com.br/nginx-ignition/core/settings"
	"dillmann.com.br/nginx-ignition/core/stream"
)

type Facade struct {
	hostCommands   *host.Commands
	streamCommands *stream.Commands
	providers      []fileProvider
	configuration  *configuration.Configuration
}

func newFacade(
	hostCommands *host.Commands,
	streamCommands *stream.Commands,
	integrationCommands *integration.Commands,
	configuration *configuration.Configuration,
	accessListRepository access_list.Repository,
	certificateRepository certificate.Repository,
	settingsRepository settings.Repository,
) *Facade {
	providers := []fileProvider{
		newAccessListFileProvider(accessListRepository),
		newHostCertificateFileProvider(certificateRepository, settingsRepository),
		newHostConfigurationFileProvider(settingsRepository, integrationCommands),
		newHostRouteSourceCodeFileProvider(),
		newMainConfigurationFileProvider(settingsRepository),
		newMimeTypesFileProvider(),
		newStreamFileProvider(),
	}

	return &Facade{
		hostCommands:   hostCommands,
		streamCommands: streamCommands,
		providers:      providers,
		configuration:  configuration,
	}
}

func (f *Facade) GetConfigurationFiles(ctx context.Context, paths *Paths) (
	configFiles []File,
	hostCount int,
	streamCount int,
	err error,
) {
	hosts, err := f.hostCommands.GetAllEnabled(ctx)
	if err != nil {
		return nil, 0, 0, err
	}

	streams, err := f.streamCommands.GetAllEnabled(ctx)
	if err != nil {
		return nil, 0, 0, err
	}

	providerCtx := &providerContext{
		context: ctx,
		paths:   paths,
		hosts:   hosts,
		streams: streams,
	}

	configFiles = make([]File, 0)
	for _, provider := range f.providers {
		files, err := provider.provide(providerCtx)
		if err != nil {
			return nil, 0, 0, err
		}

		configFiles = append(configFiles, files...)
	}

	return configFiles, len(hosts), len(streams), nil
}

func (f *Facade) ReplaceConfigurationFiles(ctx context.Context) error {
	configDir, err := f.configuration.Get("nginx-ignition.nginx.config-path")
	if err != nil {
		return err
	}

	normalizedPath := strings.TrimRight(configDir, "/")
	paths := &Paths{
		Config: normalizedPath + "/config/",
		Logs:   normalizedPath + "/logs/",
	}

	if err := f.createMissingFolders(paths); err != nil {
		return err
	}

	configFiles, hostCount, streamCount, err := f.GetConfigurationFiles(ctx, paths)
	if err != nil {
		return err
	}

	log.Infof("Rebuilding nginx configuration files for %d hosts and %d streams", hostCount, streamCount)
	if err := f.emptyConfigFolder(paths); err != nil {
		return err
	}

	for _, file := range configFiles {
		if err := f.writeConfigFile(paths, file); err != nil {
			return err
		}
	}

	return nil
}

func (f *Facade) createMissingFolders(paths *Paths) error {
	for _, folderPath := range []string{paths.Config, paths.Logs} {
		if _, err := os.Stat(folderPath); os.IsNotExist(err) {
			if err := os.MkdirAll(folderPath, os.ModePerm); err != nil {
				return fmt.Errorf("unable to create folder %s: %w", folderPath, err)
			}
		}
	}

	return nil
}

func (f *Facade) emptyConfigFolder(paths *Paths) error {
	files, err := os.ReadDir(paths.Config)
	if err != nil {
		return err
	}

	for _, file := range files {
		if err := os.RemoveAll(filepath.Join(paths.Config, file.Name())); err != nil {
			return err
		}
	}

	return nil
}

func (f *Facade) writeConfigFile(paths *Paths, file File) error {
	filePath := filepath.Join(paths.Config, file.Name)
	if err := os.WriteFile(filePath, []byte(file.Contents), os.ModePerm); err != nil {
		return fmt.Errorf("unable to write file %s: %w", filePath, err)
	}

	return nil
}
