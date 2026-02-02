package nginx

import (
	"bufio"
	"context"
	"os"
	"path/filepath"

	"dillmann.com.br/nginx-ignition/core/common/configuration"
)

type logReader struct {
	configProvider *configuration.Configuration
}

func newLogReader(configProvider *configuration.Configuration) *logReader {
	return &logReader{
		configProvider: configProvider,
	}
}

func (r *logReader) read(_ context.Context, fileName string, tailSize int) ([]LogLine, error) {
	basePath, err := r.configProvider.Get("nginx-ignition.nginx.config-path")
	if err != nil {
		return nil, err
	}

	normalizedPath := filepath.Join(basePath, "logs")
	filePath := filepath.Join(normalizedPath, fileName)
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	//nolint:errcheck
	defer file.Close()

	lines := make([]LogLine, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, LogLine{
			LineNumber: len(lines),
			Contents:   scanner.Text(),
		})
	}

	if err = scanner.Err(); err != nil {
		return nil, err
	}

	if len(lines) > tailSize {
		lines = lines[len(lines)-tailSize:]
	}

	for i := 0; i < len(lines)/2; i++ {
		lines[i], lines[len(lines)-1-i] = lines[len(lines)-1-i], lines[i]
	}

	return lines, nil
}
