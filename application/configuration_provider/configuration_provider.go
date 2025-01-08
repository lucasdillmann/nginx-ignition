package configuration_provider

import (
	"dillmann.com.br/nginx-ignition/core/common/configuration"
	"errors"
	"os"
	"strings"
)

type provider struct {
	prefix string
}

func New() configuration.Configuration {
	return &provider{}
}

func (p *provider) Get(key string) (string, error) {
	var fullKey string
	if p.prefix != "" {
		fullKey = p.prefix + "." + key
	} else {
		fullKey = key
	}

	value, exists := os.LookupEnv(fullKey)
	if exists {
		return value, nil
	}

	formattedKey := strings.ReplaceAll(key, ".", "_")
	formattedKey = strings.ReplaceAll(key, "-", "_")
	formattedKey = strings.ToUpper(formattedKey)

	value, exists = os.LookupEnv(fullKey)
	if exists {
		return value, nil
	}

	value = defaultValues()[fullKey]
	if value != "" {
		return value, nil
	}

	return "", errors.New("no configuration or environment value found for " + fullKey)
}

func (p *provider) WithPrefix(prefix string) configuration.Configuration {
	var newPrefix string
	if p.prefix == "" {
		newPrefix = prefix
	} else {
		newPrefix = p.prefix + "." + prefix
	}

	return &provider{newPrefix}
}
