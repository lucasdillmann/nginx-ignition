package integration

import (
	"context"
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
	ctx := context.Background()
	repo := New(db)

	t.Run("Save", func(t *testing.T) {
		t.Run("successfully saves a new integration", func(t *testing.T) {
			cmd := newIntegration()

			err := repo.Save(ctx, cmd)
			require.NoError(t, err)

			saved, err := repo.FindByID(ctx, cmd.ID)
			require.NoError(t, err)
			require.NotNil(t, saved)
			assert.Equal(t, cmd.Name, saved.Name)
			assert.Equal(t, cmd.Driver, saved.Driver)
			assert.Equal(t, cmd.Enabled, saved.Enabled)
			assert.Equal(t, cmd.Parameters, saved.Parameters)
		})

		t.Run("successfully updates an existing integration", func(t *testing.T) {
			id := uuid.New()
			cmd := newIntegration()
			cmd.ID = id
			require.NoError(t, repo.Save(ctx, cmd))

			cmd.Name = "Updated Integration"
			cmd.Enabled = false
			err := repo.Save(ctx, cmd)
			require.NoError(t, err)

			saved, err := repo.FindByID(ctx, id)
			require.NoError(t, err)
			assert.Equal(t, "Updated Integration", saved.Name)
			assert.False(t, saved.Enabled)
		})
	})

	t.Run("ExistsByName", func(t *testing.T) {
		t.Run("returns true when exists", func(t *testing.T) {
			cmd := newIntegration()
			require.NoError(t, repo.Save(ctx, cmd))

			exists, err := repo.ExistsByName(ctx, cmd.Name)
			require.NoError(t, err)
			assert.True(t, *exists)
		})

		t.Run("returns false when not exists", func(t *testing.T) {
			exists, err := repo.ExistsByName(ctx, "NonExistent")
			require.NoError(t, err)
			assert.False(t, *exists)
		})
	})

	t.Run("ExistsByID", func(t *testing.T) {
		t.Run("returns true when exists", func(t *testing.T) {
			cmd := newIntegration()
			require.NoError(t, repo.Save(ctx, cmd))

			exists, err := repo.ExistsByID(ctx, cmd.ID)
			require.NoError(t, err)
			assert.True(t, *exists)
		})

		t.Run("returns false when not exists", func(t *testing.T) {
			exists, err := repo.ExistsByID(ctx, uuid.New())
			require.NoError(t, err)
			assert.False(t, *exists)
		})
	})

	t.Run("FindPage", func(t *testing.T) {
		t.Run("returns a page of integrations filtered by name", func(t *testing.T) {
			prefix := uuid.New().String()
			names := []string{
				prefix + "Alpha",
				prefix + "Beta",
				prefix + "Gamma",
			}

			for _, name := range names {
				cmd := newIntegration()
				cmd.ID = uuid.New()
				cmd.Name = name
				require.NoError(t, repo.Save(ctx, cmd))
			}

			other := newIntegration()
			other.ID = uuid.New()
			other.Name = "Other" + uuid.New().String()
			require.NoError(t, repo.Save(ctx, other))

			search := prefix
			page, err := repo.FindPage(ctx, 10, 0, &search, false)
			require.NoError(t, err)

			assert.GreaterOrEqual(t, page.TotalItems, 3)

			for _, item := range page.Contents {
				assert.Contains(t, item.Name, prefix)
			}
		})

		t.Run("filters enabled only", func(t *testing.T) {
			enabled := newIntegration()
			enabled.ID = uuid.New()
			enabled.Enabled = true
			enabled.Name = uuid.New().String() + "Enabled"
			require.NoError(t, repo.Save(ctx, enabled))

			disabled := newIntegration()
			disabled.ID = uuid.New()
			disabled.Enabled = false
			disabled.Name = uuid.New().String() + "Disabled"
			require.NoError(t, repo.Save(ctx, disabled))

			search := enabled.Name
			page, err := repo.FindPage(ctx, 10, 0, &search, true)
			require.NoError(t, err)
			assert.Equal(t, 1, len(page.Contents))
			assert.Equal(t, enabled.ID, page.Contents[0].ID)
		})
	})

	t.Run("DeleteByID", func(t *testing.T) {
		t.Run("removes the integration", func(t *testing.T) {
			cmd := newIntegration()
			require.NoError(t, repo.Save(ctx, cmd))

			err := repo.DeleteByID(ctx, cmd.ID)
			require.NoError(t, err)

			exists, err := repo.ExistsByID(ctx, cmd.ID)
			require.NoError(t, err)
			assert.False(t, *exists)
		})
	})

	t.Run("InUseByID", func(t *testing.T) {
		t.Run("returns false when not in use", func(t *testing.T) {
			cmd := newIntegration()
			require.NoError(t, repo.Save(ctx, cmd))

			inUse, err := repo.InUseByID(ctx, cmd.ID)
			require.NoError(t, err)
			assert.False(t, *inUse)
		})
	})
}
