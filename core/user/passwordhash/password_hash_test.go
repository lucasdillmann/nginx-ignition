package passwordhash

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"dillmann.com.br/nginx-ignition/core/common/configuration"
)

func setupConfigForTest(t *testing.T, saltSize, iterations int) *PasswordHash {
	t.Helper()

	overrides := make(map[string]string)
	if saltSize > 0 {
		overrides["nginx-ignition.security.user-password-hashing.salt-size"] = strconv.Itoa(
			saltSize,
		)
	} else {
		overrides["nginx-ignition.security.user-password-hashing.salt-size"] = "invalid"
	}

	if iterations > 0 {
		overrides["nginx-ignition.security.user-password-hashing.iterations"] = strconv.Itoa(
			iterations,
		)
	} else {
		overrides["nginx-ignition.security.user-password-hashing.iterations"] = "invalid"
	}

	cfg := configuration.NewWithOverrides(overrides)
	passwordHash := New(cfg)

	return passwordHash
}

func Test_PasswordHash(t *testing.T) {
	t.Run("Hash", func(t *testing.T) {
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
	})

	t.Run("Verify", func(t *testing.T) {
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
	})
}
