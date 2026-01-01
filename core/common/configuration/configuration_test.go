package configuration

import (
	"fmt"
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	"dillmann.com.br/nginx-ignition/core/common/log"
)

var (
	testConfigOnce sync.Once
	testConfig     *Configuration
)

func getTestConfig() *Configuration {
	testConfigOnce.Do(func() {
		testConfig = New()
	})
	return testConfig
}

func init() {
	_ = log.Init()
}

func TestConfiguration_Get(t *testing.T) {
	cfg := getTestConfig()

	t.Run("returns value from formatted environment variable", func(t *testing.T) {
		key := "test.config.key"
		formattedKey := "TEST_CONFIG_KEY"
		expectedValue := "formatted-value"
		os.Setenv(formattedKey, expectedValue)
		defer os.Unsetenv(formattedKey)

		value, err := cfg.Get(key)

		assert.NoError(t, err)
		assert.Equal(t, expectedValue, value)
	})

	t.Run("returns value from default values", func(t *testing.T) {
		value, err := cfg.Get("nginx-ignition.server.port")

		assert.NoError(t, err)
		assert.Equal(t, "8090", value)
	})

	t.Run("returns error when value not found", func(t *testing.T) {
		_, err := cfg.Get("non-existent-key")

		assert.Error(t, err)
	})
}

func TestConfiguration_GetInt(t *testing.T) {
	cfg := getTestConfig()

	t.Run("returns integer value", func(t *testing.T) {
		key := "TEST_INT_KEY"
		os.Setenv(key, "42")
		defer os.Unsetenv(key)

		value, err := cfg.GetInt(key)

		assert.NoError(t, err)
		assert.Equal(t, 42, value)
	})

	t.Run("returns error for invalid integer", func(t *testing.T) {
		key := "TEST_INVALID_INT_KEY"
		os.Setenv(key, "not-a-number")
		defer os.Unsetenv(key)

		_, err := cfg.GetInt(key)

		assert.Error(t, err)
	})
}

func TestConfiguration_GetBoolean(t *testing.T) {
	cfg := getTestConfig()

	t.Run("returns boolean value", func(t *testing.T) {
		key := "TEST_BOOL_KEY"
		os.Setenv(key, "true")
		defer os.Unsetenv(key)

		value, err := cfg.GetBoolean(key)

		assert.NoError(t, err)
		assert.True(t, value)
	})

	t.Run("returns error for invalid boolean", func(t *testing.T) {
		key := "TEST_INVALID_BOOL_KEY"
		os.Setenv(key, "not-a-boolean")
		defer os.Unsetenv(key)

		_, err := cfg.GetBoolean(key)

		assert.Error(t, err)
	})
}

func TestConfiguration_WithPrefix(t *testing.T) {
	cfg := getTestConfig()

	t.Run("adds prefix to key lookup", func(t *testing.T) {
		key := "TEST_PREFIX_KEY"
		expectedValue := "prefixed-value"
		envKey := fmt.Sprintf("prefix.%s", key)
		os.Setenv(envKey, expectedValue)
		defer os.Unsetenv(envKey)

		prefixedCfg := cfg.WithPrefix("prefix")
		value, err := prefixedCfg.Get(key)

		assert.NoError(t, err)
		assert.Equal(t, expectedValue, value)
	})

	t.Run("chains prefixes", func(t *testing.T) {
		key := "TEST_CHAINED_KEY"
		expectedValue := "chained-value"
		envKey := fmt.Sprintf("parent.child.%s", key)
		os.Setenv(envKey, expectedValue)
		defer os.Unsetenv(envKey)

		prefixedCfg := cfg.WithPrefix("parent").WithPrefix("child")
		value, err := prefixedCfg.Get(key)

		assert.NoError(t, err)
		assert.Equal(t, expectedValue, value)
	})
}
