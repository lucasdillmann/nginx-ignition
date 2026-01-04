package stream

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
		t.Run("successfully saves a new stream", func(t *testing.T) {
			cmd := newStream()

			err := repo.Save(ctx, cmd)
			require.NoError(t, err)

			saved, err := repo.FindByID(ctx, cmd.ID)
			require.NoError(t, err)
			require.NotNil(t, saved)
			assert.Equal(t, cmd.Name, saved.Name)
			assert.Equal(t, cmd.Type, saved.Type)
			assert.Equal(t, cmd.Binding, saved.Binding)
			assert.Equal(t, cmd.DefaultBackend, saved.DefaultBackend)
			assert.Equal(t, cmd.Routes, saved.Routes)
			assert.Equal(t, cmd.FeatureSet, saved.FeatureSet)
		})

		t.Run("successfully updates an existing stream", func(t *testing.T) {
			id := uuid.New()
			cmd := newStream()
			cmd.ID = id
			require.NoError(t, repo.Save(ctx, cmd))

			cmd.Name = "Updated Stream"
			cmd.Enabled = false
			err := repo.Save(ctx, cmd)
			require.NoError(t, err)

			saved, err := repo.FindByID(ctx, id)
			require.NoError(t, err)
			assert.Equal(t, "Updated Stream", saved.Name)
			assert.False(t, saved.Enabled)
		})
	})

	t.Run("ExistsByID", func(t *testing.T) {
		t.Run("returns true when exists", func(t *testing.T) {
			cmd := newStream()
			require.NoError(t, repo.Save(ctx, cmd))

			exists, err := repo.ExistsByID(ctx, cmd.ID)
			require.NoError(t, err)
			assert.True(t, exists)
		})

		t.Run("returns false when not exists", func(t *testing.T) {
			exists, err := repo.ExistsByID(ctx, uuid.New())
			require.NoError(t, err)
			assert.False(t, exists)
		})
	})

	t.Run("FindPage", func(t *testing.T) {
		t.Run("returns a page of streams filtered by name", func(t *testing.T) {
			prefix := uuid.New().String()
			names := []string{
				prefix + "Alpha",
				prefix + "Beta",
				prefix + "Gamma",
			}

			for _, name := range names {
				cmd := newStream()
				cmd.ID = uuid.New()
				cmd.Name = name
				require.NoError(t, repo.Save(ctx, cmd))
			}

			other := newStream()
			other.ID = uuid.New()
			other.Name = "Other" + uuid.New().String()
			require.NoError(t, repo.Save(ctx, other))

			search := prefix
			page, err := repo.FindPage(ctx, 10, 0, &search)
			require.NoError(t, err)

			assert.GreaterOrEqual(t, page.TotalItems, 3)

			for _, item := range page.Contents {
				assert.Contains(t, item.Name, prefix)
			}
		})
	})

	t.Run("FindAllEnabled", func(t *testing.T) {
		t.Run("returns only enabled streams", func(t *testing.T) {
			enabled := newStream()
			enabled.ID = uuid.New()
			enabled.Enabled = true
			require.NoError(t, repo.Save(ctx, enabled))

			disabled := newStream()
			disabled.ID = uuid.New()
			disabled.Enabled = false
			require.NoError(t, repo.Save(ctx, disabled))

			all, err := repo.FindAllEnabled(ctx)
			require.NoError(t, err)

			foundEnabled := false
			foundDisabled := false
			for _, item := range all {
				if item.ID == enabled.ID {
					foundEnabled = true
				}
				if item.ID == disabled.ID {
					foundDisabled = true
				}
			}
			assert.True(t, foundEnabled)
			assert.False(t, foundDisabled)
		})
	})

	t.Run("DeleteByID", func(t *testing.T) {
		t.Run("removes the stream", func(t *testing.T) {
			cmd := newStream()
			require.NoError(t, repo.Save(ctx, cmd))

			err := repo.DeleteByID(ctx, cmd.ID)
			require.NoError(t, err)

			exists, err := repo.ExistsByID(ctx, cmd.ID)
			require.NoError(t, err)
			assert.False(t, exists)
		})
	})
}
