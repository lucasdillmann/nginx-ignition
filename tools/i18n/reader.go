package main

import (
	"bufio"
	"os"
	"strings"
)

func readPropertiesFile(path string, allKeys *map[string]bool) properties {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	props := make(properties)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			val := strings.TrimSpace(parts[1])

			props[key] = val
			(*allKeys)[key] = true
		}
	}

	return props
}
