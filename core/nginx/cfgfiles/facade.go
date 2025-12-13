package cfgfiles

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"dillmann.com.br/nginx-ignition/core/accesslist"
	"dillmann.com.br/nginx-ignition/core/certificate/server"
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
	accessListRepository accesslist.Repository,
	certificateRepository server.Repository,
	settingsRepository settings.Repository,
) *Facade {
	providers := []fileProvider{
		newAccessListFileProvider(accessListRepository),
		newHostCertificateFileProvider(certificateRepository, settingsRepository),
		newHostConfigurationFileProvider(settingsRepository, integrationCommands),
		newHostRouteStaticResponseFileProvider(),
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

func (f *Facade) GetConfigurationFiles(ctx context.Context, paths *Paths, supportedFeatures *SupportedFeatures) (
	configFiles []File,
	hosts []*host.Host,
	streams []*stream.Stream,
	err error,
) {
	enabledHosts, err := f.hostCommands.GetAllEnabled(ctx)
	if err != nil {
		return nil, nil, nil, err
	}

	enabledStreams, err := f.streamCommands.GetAllEnabled(ctx)
	if err != nil {
		return nil, nil, nil, err
	}

	providerCtx := &providerContext{
		context:           ctx,
		paths:             paths,
		hosts:             enabledHosts,
		streams:           enabledStreams,
		supportedFeatures: supportedFeatures,
	}

	configFiles = make([]File, 0)
	for _, provider := range f.providers {
		files, err := provider.provide(providerCtx)
		if err != nil {
			return nil, nil, nil, err
		}

		configFiles = append(configFiles, files...)
	}

	return configFiles, enabledHosts, enabledStreams, nil
}

func (f *Facade) ReplaceConfigurationFiles(
	ctx context.Context,
	supportedFeatures *SupportedFeatures,
) ([]*host.Host, []*stream.Stream, error) {
	configDir, err := f.configuration.Get("nginx-ignition.nginx.config-path")
	if err != nil {
		return nil, nil, err
	}

	normalizedPath := strings.TrimRight(configDir, "/")
	paths := &Paths{
		Base:   normalizedPath + "/",
		Config: normalizedPath + "/config/",
		Logs:   normalizedPath + "/logs/",
	}

	if err := f.createMissingFolders(paths); err != nil {
		return nil, nil, err
	}

	configFiles, hosts, streams, err := f.GetConfigurationFiles(ctx, paths, supportedFeatures)
	if err != nil {
		return nil, nil, err
	}

	log.Infof("Rebuilding nginx configuration files for %d hosts and %d streams", len(hosts), len(streams))
	if err := f.emptyConfigFolder(paths); err != nil {
		return nil, nil, err
	}

	for _, file := range configFiles {
		if err := f.writeConfigFile(paths, file); err != nil {
			return nil, nil, err
		}
	}

	return hosts, streams, nil
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
	if err := os.WriteFile(filePath, []byte(file.FormattedContents()), os.ModePerm); err != nil {
		return fmt.Errorf("unable to write file %s: %w", filePath, err)
	}

	return nil
}
