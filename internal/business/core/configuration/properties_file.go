package configuration

import (
	"flag"
	"os"
	"strings"
	"sync"
)

var (
	resolvedConfigFilePath string
	runOnce                sync.Once
)

func loadConfigFileValues() (map[string]string, error) {
	filePath := resolveConfigFilePath()
	file, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	output := make(map[string]string)
	for _, line := range strings.Split(string(file), "\n") {
		if strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		output[parts[0]] = parts[1]
	}

	return output, nil
}

func resolveConfigFilePath() string {
	runOnce.Do(func() {
		customPathPtr := flag.String("config", "", "Path to the configuration properties file")

		//nolint:revive
		flag.Parse()

		if customPathPtr != nil && *customPathPtr != "" {
			resolvedConfigFilePath = *customPathPtr
		} else if customPath := os.Getenv("NGINX_IGNITION_CONFIG_FILE_PATH"); customPath != "" {
			resolvedConfigFilePath = customPath
		} else {
			resolvedConfigFilePath = "nginx-ignition.properties"
		}
	})

	return resolvedConfigFilePath
}
