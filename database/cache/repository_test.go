package cache

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"dillmann.com.br/nginx-ignition/database/common/database"
	"dillmann.com.br/nginx-ignition/database/common/testutils"
)

func Test_Repository(t *testing.T) {
	testutils.RunWithMockedDatabases(t, runRepositoryTests)
}

func runRepositoryTests(t *testing.T, db *database.Database) {
	repo := New(db)

	t.Run("Save", func(t *testing.T) {
		t.Run("successfully saves a new cache configuration", func(t *testing.T) {
			cmd := newCache()

			err := repo.Save(t.Context(), cmd)
			require.NoError(t, err)

			saved, err := repo.FindByID(t.Context(), cmd.ID)
			require.NoError(t, err)
			require.NotNil(t, saved)
			assert.Equal(t, cmd.Name, saved.Name)
			assert.Equal(t, cmd.StoragePath, saved.StoragePath)
			assert.Equal(t, cmd.ConcurrencyLock, saved.ConcurrencyLock)
			assert.ElementsMatch(t, cmd.AllowedMethods, saved.AllowedMethods)
			assert.ElementsMatch(t, cmd.UseStale, saved.UseStale)
			assert.ElementsMatch(t, cmd.FileExtensions, saved.FileExtensions)
			assert.ElementsMatch(t, cmd.BypassRules, saved.BypassRules)
			assert.Len(t, saved.Durations, len(cmd.Durations))
		})

		t.Run("successfully updates an existing cache", func(t *testing.T) {
			id := uuid.New()
			cmd := newCache()
			cmd.ID = id
			cmd.Name = "Initial Name"
			cmd.MinimumUsesBeforeCaching = 1
			require.NoError(t, repo.Save(t.Context(), cmd))

			cmd.Name = "Updated Name"
			cmd.MinimumUsesBeforeCaching = 5
			err := repo.Save(t.Context(), cmd)
			require.NoError(t, err)

			saved, err := repo.FindByID(t.Context(), id)
			require.NoError(t, err)
			assert.Equal(t, "Updated Name", saved.Name)
			assert.Equal(t, 5, saved.MinimumUsesBeforeCaching)
		})
	})

	t.Run("FindByID", func(t *testing.T) {
		t.Run("returns nil when not exists", func(t *testing.T) {
			saved, err := repo.FindByID(t.Context(), uuid.New())
			require.NoError(t, err)
			assert.Nil(t, saved)
		})
	})

	t.Run("FindByID", func(t *testing.T) {
		t.Run("returns nil when not exists", func(t *testing.T) {
			saved, err := repo.FindByID(t.Context(), uuid.New())
			require.NoError(t, err)
			assert.Nil(t, saved)
		})
	})

	t.Run("FindPage", func(t *testing.T) {
		t.Run("returns a page of caches filtered by name", func(t *testing.T) {
			prefix := uuid.New().String()
			names := []string{
				prefix + "Alpha",
				prefix + "Beta",
				prefix + "Gamma",
			}

			for _, name := range names {
				cmd := newCache()
				cmd.ID = uuid.New()
				cmd.Name = name
				require.NoError(t, repo.Save(t.Context(), cmd))
			}

			other := newCache()
			other.ID = uuid.New()
			other.Name = "Other" + uuid.New().String()
			require.NoError(t, repo.Save(t.Context(), other))

			search := prefix
			page, err := repo.FindPage(t.Context(), 0, 10, &search)
			require.NoError(t, err)

			assert.GreaterOrEqual(t, page.TotalItems, 3)

			for _, item := range page.Contents {
				assert.Contains(t, item.Name, prefix)
			}
		})
	})

	t.Run("DeleteByID", func(t *testing.T) {
		t.Run("removes the cache", func(t *testing.T) {
			cmd := newCache()
			require.NoError(t, repo.Save(t.Context(), cmd))

			err := repo.DeleteByID(t.Context(), cmd.ID)
			require.NoError(t, err)

			exists, err := repo.ExistsByID(t.Context(), cmd.ID)
			require.NoError(t, err)
			assert.False(t, exists)
		})
	})

	t.Run("InUseByID", func(t *testing.T) {
		t.Run("returns false when not in use", func(t *testing.T) {
			cmd := newCache()
			require.NoError(t, repo.Save(t.Context(), cmd))

			inUse, err := repo.InUseByID(t.Context(), cmd.ID)
			require.NoError(t, err)
			assert.False(t, inUse)
		})
	})

	t.Run("FindAllInUse", func(t *testing.T) {
		t.Run("returns empty list when no caches in use", func(t *testing.T) {
			cmd := newCache()
			require.NoError(t, repo.Save(t.Context(), cmd))

			inUseList, err := repo.FindAllInUse(t.Context())
			require.NoError(t, err)
			assert.Empty(t, inUseList)
		})
	})
}
