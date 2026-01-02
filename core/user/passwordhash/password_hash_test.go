package passwordhash

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"dillmann.com.br/nginx-ignition/core/common/configuration"
	"dillmann.com.br/nginx-ignition/core/common/log"
)

func init() {
	_ = log.Init()
}

func setupConfigForTest(t *testing.T, saltSize, iterations int) *PasswordHash {
	t.Helper()
	saltKey := "NGINX_IGNITION_SECURITY_USER_PASSWORD_HASHING_SALT_SIZE"
	iterationsKey := "NGINX_IGNITION_SECURITY_USER_PASSWORD_HASHING_ITERATIONS"

	if saltSize > 0 {
		t.Setenv(saltKey, strconv.Itoa(saltSize))
	} else {
		t.Setenv(saltKey, "invalid")
	}

	if iterations > 0 {
		t.Setenv(iterationsKey, strconv.Itoa(iterations))
	} else {
		t.Setenv(iterationsKey, "invalid")
	}

	cfg := configuration.New()
	passwordHash := New(cfg)

	return passwordHash
}

func Test_PasswordHash_Hash(t *testing.T) {
	t.Run("creates a valid hash", func(t *testing.T) {
		passwordHash := setupConfigForTest(t, 32, 1)
		password := "plain-text-password"

		hash, salt, err := passwordHash.Hash(password)
		require.NoError(t, err)
		assert.NotEmpty(t, hash)
		assert.NotEmpty(t, salt)
	})

	t.Run("fails when config is invalid", func(t *testing.T) {
		passwordHash := setupConfigForTest(t, 0, 1)
		password := "plain-text-password"

		hash, salt, err := passwordHash.Hash(password)
		assert.Error(t, err)
		assert.Empty(t, hash)
		assert.Empty(t, salt)
	})
}

func Test_PasswordHash_Verify(t *testing.T) {
	t.Run("verifies a valid password", func(t *testing.T) {
		passwordHash := setupConfigForTest(t, 64, 2)
		password := "plain-text-password"

		hash, salt, err := passwordHash.Hash(password)
		require.NoError(t, err)

		ok, err := passwordHash.Verify(password, hash, salt)
		require.NoError(t, err)
		assert.True(t, ok)
	})

	t.Run("fails on wrong password", func(t *testing.T) {
		passwordHash := setupConfigForTest(t, 32, 1)
		password := "plain-text-password"

		hash, salt, err := passwordHash.Hash(password)
		require.NoError(t, err)

		ok, err := passwordHash.Verify(password+"-wrong", hash, salt)
		require.NoError(t, err)
		assert.False(t, ok)
	})

	t.Run("fails on invalid hash", func(t *testing.T) {
		passwordHash := setupConfigForTest(t, 32, 1)
		password := "plain-text-password"

		_, salt, err := passwordHash.Hash(password)
		require.NoError(t, err)

		ok, err := passwordHash.Verify(password, "invalid-hash", salt)
		assert.Error(t, err)
		assert.False(t, ok)
	})

	t.Run("fails on invalid salt", func(t *testing.T) {
		passwordHash := setupConfigForTest(t, 32, 1)
		password := "plain-text-password"

		hash, _, err := passwordHash.Hash(password)
		require.NoError(t, err)

		ok, err := passwordHash.Verify(password, hash, "invalid-salt")
		assert.Error(t, err)
		assert.False(t, ok)
	})

	t.Run("fails when config is invalid", func(t *testing.T) {
		passwordHash := setupConfigForTest(t, 32, 0)

		ok, err := passwordHash.Verify("password", "hash", "salt")
		assert.Error(t, err)
		assert.False(t, ok)
	})
}
