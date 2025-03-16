package cfgfiles

import (
	"context"
	"dillmann.com.br/nginx-ignition/core/access_list"
	"dillmann.com.br/nginx-ignition/core/certificate"
	"dillmann.com.br/nginx-ignition/core/common/configuration"
	"dillmann.com.br/nginx-ignition/core/common/log"
	"dillmann.com.br/nginx-ignition/core/host"
	"dillmann.com.br/nginx-ignition/core/integration"
	"dillmann.com.br/nginx-ignition/core/settings"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Facade struct {
	getHostsCommand host.GetAllEnabledCommand
	providers       []fileProvider
	configuration   *configuration.Configuration
}

func newFacade(
	getHostsCommand host.GetAllEnabledCommand,
	configuration *configuration.Configuration,
	accessListRepository access_list.Repository,
	certificateRepository certificate.Repository,
	settingsRepository settings.Repository,
	integrationOptionCommand integration.GetOptionUrlByIdCommand,
) *Facade {
	providers := []fileProvider{
		newAccessListFileProvider(accessListRepository),
		newHostCertificateFileProvider(certificateRepository, settingsRepository),
		newHostConfigurationFileProvider(settingsRepository, integrationOptionCommand),
		newHostRouteSourceCodeFileProvider(),
		newMainConfigurationFileProvider(settingsRepository),
		newMimeTypesFileProvider(),
	}

	return &Facade{
		getHostsCommand: getHostsCommand,
		providers:       providers,
		configuration:   configuration,
	}
}

func (f *Facade) ReplaceConfigurationFiles(ctx context.Context) error {
	hosts, err := f.getHostsCommand(ctx)
	if err != nil {
		return err
	}

	log.Infof("Rebuilding nginx configuration files for %d hosts", len(hosts))

	configDir, err := f.configuration.Get("nginx-ignition.nginx.config-path")
	if err != nil {
		return err
	}

	normalizedPath := strings.TrimRight(configDir, "/")
	if err := f.createMissingFolders(normalizedPath); err != nil {
		return err
	}

	if err := f.emptyConfigFolder(normalizedPath); err != nil {
		return err
	}

	var configFiles []output
	for _, provider := range f.providers {
		files, err := provider.provide(ctx, normalizedPath, hosts)
		if err != nil {
			return err
		}

		configFiles = append(configFiles, files...)
	}

	for _, file := range configFiles {
		if err := f.writeConfigFile(normalizedPath, file); err != nil {
			return err
		}
	}

	return nil
}

func (f *Facade) createMissingFolders(basePath string) error {
	folders := []string{"logs", "config"}
	for _, folderName := range folders {
		folderPath := filepath.Join(basePath, folderName)
		if _, err := os.Stat(folderPath); os.IsNotExist(err) {
			if err := os.MkdirAll(folderPath, os.ModePerm); err != nil {
				return fmt.Errorf("unable to create folder %s: %w", folderPath, err)
			}
		}
	}
	return nil
}

func (f *Facade) emptyConfigFolder(basePath string) error {
	configPath := filepath.Join(basePath, "config")
	files, err := os.ReadDir(configPath)
	if err != nil {
		return err
	}
	for _, file := range files {
		if err := os.RemoveAll(filepath.Join(configPath, file.Name())); err != nil {
			return err
		}
	}
	return nil
}

func (f *Facade) writeConfigFile(basePath string, file output) error {
	filePath := filepath.Join(basePath, "config", file.name)
	if err := os.WriteFile(filePath, []byte(file.contents), os.ModePerm); err != nil {
		return fmt.Errorf("unable to write file %s: %w", filePath, err)
	}
	return nil
}
