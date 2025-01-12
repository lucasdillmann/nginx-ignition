package configuration

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Configuration struct {
	prefix string
}

func New() *Configuration {
	return &Configuration{}
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

	value = defaultValues()[fullKey]
	if value != "" {
		return value, nil
	}

	return "", errors.New(fmt.Sprintf(
		"no configuration or environment value found for %s or %s",
		fullKey,
		formattedKey,
	))
}

func (c *Configuration) GetInt(key string) (int, error) {
	value, err := c.Get(key)
	if err != nil {
		return 0, err
	}

	return strconv.Atoi(value)
}

func (c *Configuration) WithPrefix(prefix string) *Configuration {
	var newPrefix string
	if c.prefix == "" {
		newPrefix = prefix
	} else {
		newPrefix = c.prefix + "." + prefix
	}

	return &Configuration{newPrefix}
}
