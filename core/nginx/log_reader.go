package nginx

import (
	"bufio"
	"dillmann.com.br/nginx-ignition/core/common/configuration"
	"os"
	"path/filepath"
	"strings"
)

type logReader struct {
	configProvider *configuration.Configuration
}

func newLogReader(configProvider *configuration.Configuration) *logReader {
	return &logReader{
		configProvider: configProvider,
	}
}

func (r *logReader) read(fileName string, tailSize int) ([]string, error) {
	basePath, err := r.configProvider.Get("nginx-ignition.nginx.config-directory")
	if err != nil {
		return nil, err
	}

	normalizedPath := strings.TrimRight(basePath, "/") + "/logs"

	filePath := filepath.Join(normalizedPath, fileName)
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err = scanner.Err(); err != nil {
		return nil, err
	}

	if len(lines) > tailSize {
		lines = lines[len(lines)-tailSize:]
	}

	return lines, nil
}
