package configuration

import (
	"os"
	"strings"
)

func loadConfigFileValues() (map[string]string, error) {
	customPath := os.Getenv("NGINX_IGNITION_CONFIG_FILE_PATH")
	if customPath == "" {
		customPath = "nginx-ignition.properties"
	}

	file, err := os.ReadFile(customPath)
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
