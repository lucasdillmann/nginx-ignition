package configuration

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"dillmann.com.br/nginx-ignition/core/common/log"
)

type Configuration struct {
	configFileValues map[string]string
	prefix           string
}

func New() *Configuration {
	configFileValues, err := loadConfigFileValues()
	if err != nil {
		log.Warnf("Unable to read configuration properties file: %s", err)
	}

	return &Configuration{
		configFileValues: configFileValues,
	}
}

func NewWithOverrides(overrides map[string]string) *Configuration {
	values, err := loadConfigFileValues()
	if err != nil {
		log.Warnf("Unable to read configuration properties file: %s", err)
		values = make(map[string]string)
	}

	for key, value := range overrides {
		values[key] = value
	}

	return &Configuration{
		configFileValues: values,
	}
}

func (c *Configuration) Get(key string) (string, error) {
	var fullKey string
	if c.prefix != "" {
		fullKey = c.prefix + "." + key
	} else {
		fullKey = key
	}

	value, exists := os.LookupEnv(fullKey)
	if exists {
		return value, nil
	}

	formattedKey := strings.ReplaceAll(fullKey, ".", "_")
	formattedKey = strings.ReplaceAll(formattedKey, "-", "_")
	formattedKey = strings.ToUpper(formattedKey)

	value, exists = os.LookupEnv(formattedKey)
	if exists {
		return value, nil
	}

	value = c.configFileValues[fullKey]
	if value != "" {
		return value, nil
	}

	value = defaultValues[fullKey]
	if value != "" {
		return value, nil
	}

	return "", fmt.Errorf(
		"no configuration or environment value found for %s or %s",
		fullKey,
		formattedKey,
	)
}

func (c *Configuration) GetInt(key string) (int, error) {
	value, err := c.Get(key)
	if err != nil {
		return 0, err
	}

	return strconv.Atoi(value)
}

func (c *Configuration) GetBoolean(key string) (bool, error) {
	value, err := c.Get(key)
	if err != nil {
		return false, err
	}

	return strconv.ParseBool(value)
}

func (c *Configuration) WithPrefix(prefix string) *Configuration {
	var newPrefix string
	if c.prefix == "" {
		newPrefix = prefix
	} else {
		newPrefix = c.prefix + "." + prefix
	}

	return &Configuration{
		configFileValues: c.configFileValues,
		prefix:           newPrefix,
	}
}
