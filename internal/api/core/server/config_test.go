package server

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lucasdillmann/nginx-ignition/internal/business/core/configuration"
)

var serverConfigFields = []struct {
	key   string
	value string
}{
	{key: "port", value: "8090"},
	{key: "address", value: "0.0.0.0"},
	{key: "read-timeout-seconds", value: "15"},
	{key: "write-timeout-seconds", value: "30"},
	{key: "idle-timeout-seconds", value: "120"},
	{key: "read-header-timeout-seconds", value: "2"},
	{key: "max-header-bytes", value: "16384"},
}

func Test_loadServerConfig(t *testing.T) {
	t.Run("loads defaults", func(t *testing.T) {
		clearServerConfigEnvironment(t)

		actual, err := loadServerConfig(configuration.NewWithOverrides(map[string]string{}))

		require.NoError(t, err)
		assert.Equal(t, &serverConfig{
			port:              "8090",
			address:           "0.0.0.0",
			readTimeout:       15 * time.Second,
			writeTimeout:      30 * time.Second,
			idleTimeout:       120 * time.Second,
			readHeaderTimeout: 2 * time.Second,
			maxHeaderBytes:    16384,
		}, actual)
	})

	t.Run("loads all custom values and converts durations", func(t *testing.T) {
		clearServerConfigEnvironment(t)
		overrides := map[string]string{
			"nginx-ignition.server.port":                        "9191",
			"nginx-ignition.server.address":                     "127.0.0.1",
			"nginx-ignition.server.read-timeout-seconds":        "11",
			"nginx-ignition.server.write-timeout-seconds":       "22",
			"nginx-ignition.server.idle-timeout-seconds":        "33",
			"nginx-ignition.server.read-header-timeout-seconds": "44",
			"nginx-ignition.server.max-header-bytes":            "65536",
		}

		actual, err := loadServerConfig(configuration.NewWithOverrides(overrides))

		require.NoError(t, err)
		assert.Equal(t, &serverConfig{
			port:              "9191",
			address:           "127.0.0.1",
			readTimeout:       11 * time.Second,
			writeTimeout:      22 * time.Second,
			idleTimeout:       33 * time.Second,
			readHeaderTimeout: 44 * time.Second,
			maxHeaderBytes:    65536,
		}, actual)
	})

	t.Run("uses environment values with higher precedence", func(t *testing.T) {
		clearServerConfigEnvironment(t)
		t.Setenv("NGINX_IGNITION_SERVER_PORT", "9292")
		overrides := map[string]string{
			"nginx-ignition.server.port": "9191",
		}

		actual, err := loadServerConfig(configuration.NewWithOverrides(overrides))

		require.NoError(t, err)
		assert.Equal(t, "9292", actual.port)
	})

	t.Run("returns nil and an error for every missing value", func(t *testing.T) {
		clearServerConfigEnvironment(t)
		tests := []struct {
			key  string
			name string
		}{
			{name: "port", key: "port"},
			{name: "address", key: "address"},
			{name: "read timeout", key: "read-timeout-seconds"},
			{name: "write timeout", key: "write-timeout-seconds"},
			{name: "idle timeout", key: "idle-timeout-seconds"},
			{name: "read header timeout", key: "read-header-timeout-seconds"},
			{name: "max header bytes", key: "max-header-bytes"},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				prefix := "load-server-config-missing"
				overrides := completeServerConfigOverrides(prefix)
				delete(overrides, prefix+".nginx-ignition.server."+test.key)
				cfg := configuration.NewWithOverrides(overrides).WithPrefix(prefix)

				actual, err := loadServerConfig(cfg)

				assert.Nil(t, actual)
				assert.ErrorContains(t, err, prefix+".nginx-ignition.server."+test.key)
			})
		}
	})

	t.Run("returns nil and an error for every invalid integer", func(t *testing.T) {
		clearServerConfigEnvironment(t)
		tests := []struct {
			key  string
			name string
		}{
			{name: "read timeout", key: "read-timeout-seconds"},
			{name: "write timeout", key: "write-timeout-seconds"},
			{name: "idle timeout", key: "idle-timeout-seconds"},
			{name: "read header timeout", key: "read-header-timeout-seconds"},
			{name: "max header bytes", key: "max-header-bytes"},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				prefix := "load-server-config-invalid"
				overrides := completeServerConfigOverrides(prefix)
				overrides[prefix+".nginx-ignition.server."+test.key] = "invalid"
				cfg := configuration.NewWithOverrides(overrides).WithPrefix(prefix)

				actual, err := loadServerConfig(cfg)

				assert.Nil(t, actual)
				assert.Error(t, err)
			})
		}
	})

	t.Run("panics for a nil configuration", func(t *testing.T) {
		assert.Panics(t, func() {
			_, _ = loadServerConfig(nil)
		})
	})
}

func clearServerConfigEnvironment(t *testing.T) {
	t.Helper()

	prefix := "nginx-ignition.server."
	for _, field := range serverConfigFields {
		keys := []string{
			prefix + field.key,
			"NGINX_IGNITION_SERVER_" + field.key,
		}
		keys[1] = "NGINX_IGNITION_SERVER_"
		for _, character := range field.key {
			switch character {
			case '-':
				keys[1] += "_"
			default:
				keys[1] += string(character)
			}
		}

		for _, key := range keys {
			value, exists := os.LookupEnv(key)
			require.NoError(t, os.Unsetenv(key))
			if exists {
				t.Cleanup(func() {
					require.NoError(t, os.Setenv(key, value))
				})
			}
		}
	}
}

func completeServerConfigOverrides(prefix string) map[string]string {
	values := make(map[string]string, len(serverConfigFields))
	for _, field := range serverConfigFields {
		key := "nginx-ignition.server." + field.key
		if prefix != "" {
			key = prefix + "." + key
		}

		values[key] = field.value
	}

	return values
}
