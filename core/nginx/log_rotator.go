package nginx

import (
	"bufio"
	"context"
	"dillmann.com.br/nginx-ignition/core/common/configuration"
	"dillmann.com.br/nginx-ignition/core/common/log"
	"dillmann.com.br/nginx-ignition/core/host"
	"dillmann.com.br/nginx-ignition/core/settings"
	"os"
	"path/filepath"
	"strings"
)

type logRotator struct {
	configProvider     *configuration.Configuration
	settingsRepository settings.Repository
	hostRepository     host.Repository
	processManager     *processManager
}

func newLogRotator(
	configProvider *configuration.Configuration,
	settingsRepository settings.Repository,
	hostRepository host.Repository,
	processManager *processManager,
) *logRotator {
	return &logRotator{
		configProvider:     configProvider,
		settingsRepository: settingsRepository,
		hostRepository:     hostRepository,
		processManager:     processManager,
	}
}

func (r *logRotator) rotate(ctx context.Context) error {
	log.Infof("Starting log rotation")

	basePath, err := r.configProvider.Get("nginx-ignition.nginx.config-path")
	if err != nil {
		return err
	}

	normalizedPath := strings.TrimRight(basePath, "/") + "/logs"

	cfg, err := r.settingsRepository.Get(ctx)
	if err != nil {
		return err
	}

	maximumLines := cfg.LogRotation.MaximumLines

	logFiles, err := r.getLogFiles(ctx)
	if err != nil {
		return err
	}

	for _, logFile := range logFiles {
		if err = r.rotateFile(ctx, normalizedPath, logFile, maximumLines); err != nil {
			log.Warnf("Unable to rotate log file %s: %v", logFile, err)
		}
	}

	if err = r.processManager.sendReopenSignal(); err != nil {
		return err
	}

	log.Infof("Log rotation finished with %d files trimmed to up to %d lines", len(logFiles), maximumLines)
	return nil
}

func (r *logRotator) rotateFile(_ context.Context, basePath, fileName string, maximumLines int) error {
	filePath := filepath.Join(basePath, fileName)
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}

	defer file.Close()

	lines, err := r.readTail(file, maximumLines)
	if err != nil {
		return err
	}

	if len(lines) < maximumLines {
		return nil
	}

	trimmedContent := strings.Join(lines, "\n") + "\n"
	return os.WriteFile(filePath, []byte(trimmedContent), 0777)
}

func (r *logRotator) readTail(file *os.File, size int) ([]string, error) {
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if len(lines) > size {
		lines = lines[len(lines)-size:]
	}

	return lines, nil
}

func (r *logRotator) getLogFiles(ctx context.Context) ([]string, error) {
	hosts, err := r.hostRepository.FindAllEnabled(ctx)
	if err != nil {
		return nil, err
	}

	var logFiles []string
	for _, item := range hosts {
		logFiles = append(
			logFiles,
			"host-"+item.ID.String()+".access.log",
			"host-"+item.ID.String()+".error.log",
		)
	}

	logFiles = append(logFiles, "main.log")

	return logFiles, nil
}
