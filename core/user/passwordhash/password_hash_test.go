package passwordhash_test

import (
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"dillmann.com.br/nginx-ignition/core/common/configuration"
	"dillmann.com.br/nginx-ignition/core/common/log"
	"dillmann.com.br/nginx-ignition/core/user/passwordhash"
)

func init() {
	log.Init()
}

func setupConfigForTest(t *testing.T, saltSize, iterations int) (*configuration.Configuration, func()) {
	t.Helper()

	saltKey := "NGINX_IGNITION_SECURITY_USER_PASSWORD_HASHING_SALT_SIZE"
	iterationsKey := "NGINX_IGNITION_SECURITY_USER_PASSWORD_HASHING_ITERATIONS"

	os.Setenv(saltKey, strconv.Itoa(saltSize))
	os.Setenv(iterationsKey, strconv.Itoa(iterations))

	cfg := configuration.New()

	cleanup := func() {
		os.Unsetenv(saltKey)
		os.Unsetenv(iterationsKey)
	}

	return cfg, cleanup
}

func TestPasswordHash_HashAndVerify(t *testing.T) {
	t.Run("valid password", func(t *testing.T) {
		cfg, cleanup := setupConfigForTest(t, 32, 1)
		defer cleanup()

		h := passwordhash.New(cfg)
		password := "plain-text-password"

		hash, salt, err := h.Hash(password)
		require.NoError(t, err)

		ok, err := h.Verify(password, hash, salt)
		require.NoError(t, err)
		assert.True(t, ok)
	})

	t.Run("wrong password", func(t *testing.T) {
		cfg, cleanup := setupConfigForTest(t, 32, 1)
		defer cleanup()

		h := passwordhash.New(cfg)
		password := "plain-text-password"

		hash, salt, err := h.Hash(password)
		require.NoError(t, err)

		ok, err := h.Verify(password+"-wrong", hash, salt)
		require.NoError(t, err)
		assert.False(t, ok)
	})

	t.Run("invalid hash", func(t *testing.T) {
		cfg, cleanup := setupConfigForTest(t, 32, 1)
		defer cleanup()

		h := passwordhash.New(cfg)
		password := "plain-text-password"

		_, salt, err := h.Hash(password)
		require.NoError(t, err)

		ok, err := h.Verify(password, "invalid-hash", salt)
		assert.Error(t, err)
		assert.False(t, ok)
	})

	t.Run("invalid salt", func(t *testing.T) {
		cfg, cleanup := setupConfigForTest(t, 32, 1)
		defer cleanup()

		h := passwordhash.New(cfg)
		password := "plain-text-password"

		hash, _, err := h.Hash(password)
		require.NoError(t, err)

		ok, err := h.Verify(password, hash, "invalid-salt")
		assert.Error(t, err)
		assert.False(t, ok)
	})
}

func TestPasswordHash_ErrorOnMissingConfig(t *testing.T) {
	saltKey := "NGINX_IGNITION_SECURITY_USER_PASSWORD_HASHING_SALT_SIZE"
	iterationsKey := "NGINX_IGNITION_SECURITY_USER_PASSWORD_HASHING_ITERATIONS"

	os.Unsetenv(saltKey)
	os.Unsetenv(iterationsKey)

	h := passwordhash.New(configuration.New())
	password := "plain-text-password"

	hash, salt, err := h.Hash(password)
	assert.Error(t, err)
	assert.Empty(t, hash)
	assert.Empty(t, salt)

	ok, err := h.Verify(password, "hash", "salt")
	assert.Error(t, err)
	assert.False(t, ok)
}
