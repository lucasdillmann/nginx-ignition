package configuration

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Configuration(t *testing.T) {
	t.Run("Get", func(t *testing.T) {
		t.Run("returns value from formatted environment variable", func(t *testing.T) {
			config := newConfiguration()
			key := "test.config.key"
			formattedKey := "TEST_CONFIG_KEY"
			expectedValue := "formatted-value"
			t.Setenv(formattedKey, expectedValue)

			value, err := config.Get(key)

			assert.NoError(t, err)
			assert.Equal(t, expectedValue, value)
		})

		t.Run("returns value from default values", func(t *testing.T) {
			config := newConfiguration()
			value, err := config.Get("nginx-ignition.server.port")

			assert.NoError(t, err)
			assert.Equal(t, "8090", value)
		})

		t.Run("returns error when value not found", func(t *testing.T) {
			config := newConfiguration()
			_, err := config.Get("non-existent-key")

			assert.Error(t, err)
		})
	})

	t.Run("GetInt", func(t *testing.T) {
		t.Run("returns integer value", func(t *testing.T) {
			config := newConfiguration()
			key := "TEST_INT_KEY"
			t.Setenv(key, "42")

			value, err := config.GetInt(key)

			assert.NoError(t, err)
			assert.Equal(t, 42, value)
		})

		t.Run("returns error for invalid integer", func(t *testing.T) {
			config := newConfiguration()
			key := "TEST_INVALID_INT_KEY"
			t.Setenv(key, "not-a-number")

			_, err := config.GetInt(key)

			assert.Error(t, err)
		})
	})

	t.Run("GetBoolean", func(t *testing.T) {
		t.Run("returns boolean value", func(t *testing.T) {
			config := newConfiguration()
			key := "TEST_BOOL_KEY"
			t.Setenv(key, "true")

			value, err := config.GetBoolean(key)

			assert.NoError(t, err)
			assert.True(t, value)
		})

		t.Run("returns error for invalid boolean", func(t *testing.T) {
			config := newConfiguration()
			key := "TEST_INVALID_BOOL_KEY"
			t.Setenv(key, "not-a-boolean")

			_, err := config.GetBoolean(key)

			assert.Error(t, err)
		})
	})

	t.Run("WithPrefix", func(t *testing.T) {
		t.Run("adds prefix to key lookup", func(t *testing.T) {
			config := newConfiguration()
			key := "TEST_PREFIX_KEY"
			expectedValue := "prefixed-value"
			envKey := fmt.Sprintf("prefix.%s", key)
			t.Setenv(envKey, expectedValue)

			prefixedConfig := config.WithPrefix("prefix")
			value, err := prefixedConfig.Get(key)

			assert.NoError(t, err)
			assert.Equal(t, expectedValue, value)
		})

		t.Run("chains prefixes", func(t *testing.T) {
			config := newConfiguration()
			key := "TEST_CHAINED_KEY"
			expectedValue := "chained-value"
			envKey := fmt.Sprintf("parent.child.%s", key)
			t.Setenv(envKey, expectedValue)

			prefixedConfig := config.WithPrefix("parent").WithPrefix("child")
			value, err := prefixedConfig.Get(key)

			assert.NoError(t, err)
			assert.Equal(t, expectedValue, value)
		})
	})
}
