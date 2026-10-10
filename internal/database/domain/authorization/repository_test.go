package authorization

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lucasdillmann/nginx-ignition/internal/database/core/database"
	"github.com/lucasdillmann/nginx-ignition/internal/database/core/testutils"
)

const (
	testJwtSecret        = "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
	testAnotherJwtSecret = "fedcba9876543210fedcba9876543210fedcba9876543210fedcba9876543210"
)

func Test_Repository(t *testing.T) {
	testutils.RunWithMockedDatabases(t, runRepositoryTests)
}

func runRepositoryTests(t *testing.T, db *database.Database) {
	repo := New(db)

	t.Run("FindJwtSecret", func(t *testing.T) {
		t.Run("returns nil when the table is empty", func(t *testing.T) {
			found, err := repo.FindJwtSecret(t.Context())
			require.NoError(t, err)
			assert.Nil(t, found)
		})

		t.Run("returns the stored JWT secret", func(t *testing.T) {
			secret := testJwtSecret
			stored, err := repo.SaveJwtSecretIfNotExists(t.Context(), &secret)
			require.NoError(t, err)
			require.NotNil(t, stored)

			found, err := repo.FindJwtSecret(t.Context())
			require.NoError(t, err)
			require.NotNil(t, found)
			assert.Equal(t, testJwtSecret, *found)
		})
	})

	t.Run("SaveJwtSecretIfNotExists", func(t *testing.T) {
		t.Run("stores the JWT secret when the table is empty", func(t *testing.T) {
			secret := testJwtSecret
			stored, err := repo.SaveJwtSecretIfNotExists(t.Context(), &secret)
			require.NoError(t, err)
			require.NotNil(t, stored)
			assert.Equal(t, testJwtSecret, *stored)

			total, err := db.Select().Model((*configurationModel)(nil)).Count(t.Context())
			require.NoError(t, err)
			assert.Equal(t, int64(1), total)
		})

		t.Run(
			"keeps the already stored JWT secret and returns it instead of the provided one",
			func(t *testing.T) {
				secret := testJwtSecret
				_, err := repo.SaveJwtSecretIfNotExists(t.Context(), &secret)
				require.NoError(t, err)

				another := testAnotherJwtSecret
				stored, err := repo.SaveJwtSecretIfNotExists(t.Context(), &another)
				require.NoError(t, err)
				require.NotNil(t, stored)
				assert.Equal(t, testJwtSecret, *stored)

				found, err := repo.FindJwtSecret(t.Context())
				require.NoError(t, err)
				require.NotNil(t, found)
				assert.Equal(t, testJwtSecret, *found)

				total, err := db.Select().Model((*configurationModel)(nil)).Count(t.Context())
				require.NoError(t, err)
				assert.Equal(t, int64(1), total)
			},
		)
	})
}
