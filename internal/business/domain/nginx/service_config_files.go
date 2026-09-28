package nginx

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/lucasdillmann/nginx-ignition/internal/business/core/log"
	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/host"
	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/nginx/cfgfiles/provider"
)

func (s *service) GetConfigFiles(
	ctx context.Context,
	input GetConfigFilesInput,
) ([]byte, error) {
	paths := &provider.Paths{
		Base:   input.BasePath,
		Config: input.ConfigPath,
		Logs:   input.LogPath,
		Cache:  input.CachePath,
		Temp:   input.TempPath,
	}

	supportedFeatures, err := s.resolveSupportedFeatures(ctx)
	if err != nil {
		return nil, err
	}

	configFiles, _, _, err := s.configFilesManager.GetConfigurationFiles(
		ctx,
		paths,
		supportedFeatures,
	)
	if err != nil {
		return nil, err
	}

	buffer := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buffer)

	//nolint:errcheck
	defer zipWriter.Close()

	for _, file := range configFiles {
		itemWriter, err := zipWriter.Create(file.Name)
		if err != nil {
			return nil, err
		}

		if _, err := itemWriter.Write([]byte(file.FormattedContents())); err != nil {
			return nil, err
		}
	}

	if err := zipWriter.Close(); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (s *service) replaceConfigurationFiles(
	ctx context.Context,
	supportedFeatures *provider.SupportedFeatures,
) ([]host.Host, error) {
	paths := s.configPaths()
	if err := s.createMissingFolders(paths); err != nil {
		return nil, err
	}

	log.Infof("Rebuilding nginx configuration files on %s...", paths.Base)
	if err := s.testConfigurationFiles(ctx, supportedFeatures, paths); err != nil {
		return nil, err
	}

	configFiles, hosts, _, err := s.configFilesManager.GetConfigurationFiles(
		ctx,
		paths,
		supportedFeatures,
	)
	if err != nil {
		return nil, err
	}

	log.Infof(
		"Validation and generation successful. Creating (or replacing) configuration files...",
	)
	if err := s.emptyConfigFolder(paths); err != nil {
		return nil, err
	}

	for _, file := range configFiles {
		if err := s.writeConfigFile(paths, file); err != nil {
			return nil, err
		}
	}

	return hosts, nil
}

func (s *service) testConfigurationFiles(
	ctx context.Context,
	supportedFeatures *provider.SupportedFeatures,
	paths *provider.Paths,
) error {
	testRoot, err := s.createTestConfigRoot(paths)
	if err != nil {
		return err
	}

	defer s.removeTestConfigRoot(testRoot)

	testPaths := newConfigPaths(testRoot)
	if err = s.createMissingFolders(testPaths); err != nil {
		return err
	}

	configFiles, _, _, err := s.configFilesManager.GetConfigurationFiles(
		ctx,
		testPaths,
		supportedFeatures,
	)
	if err != nil {
		return err
	}

	for _, file := range configFiles {
		if err := s.writeConfigFile(testPaths, file); err != nil {
			return err
		}
	}

	log.Infof("Testing nginx staging configuration files on %s...", testRoot)
	return s.processManager.ValidateConfiguration(ctx, testRoot)
}

func (s *service) createTestConfigRoot(paths *provider.Paths) (string, error) {
	root, err := os.MkdirTemp(filepath.Clean(paths.Temp), "config-test-")
	if err != nil {
		return "", fmt.Errorf("unable to create configuration test folder: %w", err)
	}

	return root, nil
}

func (s *service) removeTestConfigRoot(root string) {
	if err := os.RemoveAll(root); err != nil {
		log.Warnf("Unable to remove configuration test folder %s: %v", root, err)
	}
}

func (s *service) configPaths() *provider.Paths {
	return newConfigPaths(s.processManager.configPath)
}

func newConfigPaths(root string) *provider.Paths {
	cleanPath := filepath.Clean(root)
	toNginxPath := func(p string) string {
		return filepath.ToSlash(p) + "/"
	}

	return &provider.Paths{
		Base:   toNginxPath(cleanPath),
		Config: toNginxPath(filepath.Join(cleanPath, "config")),
		Logs:   toNginxPath(filepath.Join(cleanPath, "logs")),
		Cache:  toNginxPath(filepath.Join(cleanPath, "cache")),
		Temp:   toNginxPath(filepath.Join(cleanPath, "temp")),
	}
}

func (s *service) createMissingFolders(paths *provider.Paths) error {
	for _, folderPath := range []string{paths.Config, paths.Logs, paths.Cache, paths.Temp} {
		if _, err := os.Stat(folderPath); os.IsNotExist(err) {
			if err := os.MkdirAll(folderPath, os.ModePerm); err != nil {
				return fmt.Errorf("unable to create folder %s: %w", folderPath, err)
			}
		}
	}

	return nil
}

func (s *service) emptyConfigFolder(paths *provider.Paths) error {
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

func (s *service) writeConfigFile(paths *provider.Paths, file provider.File) error {
	filePath := filepath.Join(paths.Config, file.Name)
	if err := os.WriteFile(filePath, []byte(file.FormattedContents()), 0o644); err != nil {
		return fmt.Errorf("unable to write file %s: %w", filePath, err)
	}

	return nil
}
