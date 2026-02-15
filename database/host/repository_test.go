package host

import (
	"strings"
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
		t.Run("successfully saves a new host", func(t *testing.T) {
			cmd := newHost()

			err := repo.Save(t.Context(), cmd)
			require.NoError(t, err)

			saved, err := repo.FindByID(t.Context(), cmd.ID)
			require.NoError(t, err)
			require.NotNil(t, saved)
			assert.Equal(t, cmd.DomainNames, saved.DomainNames)
			assert.Equal(t, cmd.Enabled, saved.Enabled)
			assert.Equal(t, cmd.DefaultServer, saved.DefaultServer)
			assert.Equal(t, cmd.UseGlobalBindings, saved.UseGlobalBindings)
			assert.Equal(t, cmd.FeatureSet, saved.FeatureSet)
			assert.Len(t, saved.Bindings, len(cmd.Bindings))
			assert.Len(t, saved.Routes, len(cmd.Routes))
		})

		t.Run("successfully updates an existing host", func(t *testing.T) {
			id := uuid.New()
			cmd := newHost()
			cmd.ID = id
			require.NoError(t, repo.Save(t.Context(), cmd))

			cmd.Enabled = false
			cmd.DomainNames = []string{"updated.example.com"}
			err := repo.Save(t.Context(), cmd)
			require.NoError(t, err)

			saved, err := repo.FindByID(t.Context(), id)
			require.NoError(t, err)
			assert.ElementsMatch(t, []string{"updated.example.com"}, saved.DomainNames)
			assert.False(t, saved.Enabled)
		})
	})

	t.Run("ExistsByID", func(t *testing.T) {
		t.Run("returns true when exists", func(t *testing.T) {
			cmd := newHost()
			require.NoError(t, repo.Save(t.Context(), cmd))

			exists, err := repo.ExistsByID(t.Context(), cmd.ID)
			require.NoError(t, err)
			assert.True(t, exists)
		})

		t.Run("returns false when not exists", func(t *testing.T) {
			exists, err := repo.ExistsByID(t.Context(), uuid.New())
			require.NoError(t, err)
			assert.False(t, exists)
		})
	})

	t.Run("FindPage", func(t *testing.T) {
		t.Run("returns a page of hosts filtered by domain name", func(t *testing.T) {
			prefix := uuid.New().String()
			domains := []string{
				prefix + ".com",
				prefix + ".org",
				prefix + ".net",
			}

			for _, domain := range domains {
				cmd := newHost()
				cmd.ID = uuid.New()
				cmd.DomainNames = []string{domain}
				require.NoError(t, repo.Save(t.Context(), cmd))
			}

			other := newHost()
			other.ID = uuid.New()
			other.DomainNames = []string{"other.com"}
			require.NoError(t, repo.Save(t.Context(), other))

			page, err := repo.FindPage(t.Context(), 10, 0, new(prefix))
			require.NoError(t, err)

			assert.GreaterOrEqual(t, page.TotalItems, 3)

			for _, item := range page.Contents {
				found := false
				for _, d := range item.DomainNames {
					if strings.Contains(d, prefix) {
						found = true
						break
					}
				}
				assert.True(t, found)
			}
		})
	})

	t.Run("FindAllEnabled", func(t *testing.T) {
		t.Run("returns only enabled hosts", func(t *testing.T) {
			enabled := newHost()
			enabled.ID = uuid.New()
			enabled.Enabled = true
			require.NoError(t, repo.Save(t.Context(), enabled))

			disabled := newHost()
			disabled.ID = uuid.New()
			disabled.Enabled = false
			require.NoError(t, repo.Save(t.Context(), disabled))

			all, err := repo.FindAllEnabled(t.Context())
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

	t.Run("FindDefault", func(t *testing.T) {
		t.Run("returns the default server", func(t *testing.T) {
			cleanup(t.Context(), t, repo)

			def := newHost()
			def.ID = uuid.New()
			def.DefaultServer = true
			def.DomainNames = nil
			require.NoError(t, repo.Save(t.Context(), def))

			normal := newHost()
			normal.ID = uuid.New()
			normal.DefaultServer = false
			require.NoError(t, repo.Save(t.Context(), normal))

			found, err := repo.FindDefault(t.Context())
			require.NoError(t, err)
			require.NotNil(t, found)
			assert.Equal(t, def.ID, found.ID)
		})

		t.Run("returns nil when no default server", func(t *testing.T) {
			cleanup(t.Context(), t, repo)

			existing, err := repo.FindDefault(t.Context())
			require.NoError(t, err)
			if existing != nil {
				err = repo.DeleteByID(t.Context(), existing.ID)
				require.NoError(t, err)
			}

			found, err := repo.FindDefault(t.Context())
			require.NoError(t, err)
			assert.Nil(t, found)
		})
	})

	t.Run("DeleteByID", func(t *testing.T) {
		t.Run("removes the host", func(t *testing.T) {
			cmd := newHost()
			require.NoError(t, repo.Save(t.Context(), cmd))

			err := repo.DeleteByID(t.Context(), cmd.ID)
			require.NoError(t, err)

			exists, err := repo.ExistsByID(t.Context(), cmd.ID)
			require.NoError(t, err)
			assert.False(t, exists)
		})
	})
}
