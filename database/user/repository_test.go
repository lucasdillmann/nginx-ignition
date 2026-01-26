package user

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
		t.Run("successfully saves a new user", func(t *testing.T) {
			cmd := newUser()

			err := repo.Save(t.Context(), cmd)
			require.NoError(t, err)

			saved, err := repo.FindByID(t.Context(), cmd.ID)
			require.NoError(t, err)
			require.NotNil(t, saved)
			assert.Equal(t, cmd.Name, saved.Name)
			assert.Equal(t, cmd.Username, saved.Username)
			assert.Equal(t, cmd.PasswordHash, saved.PasswordHash)
			assert.Equal(t, cmd.PasswordSalt, saved.PasswordSalt)
			assert.Equal(t, cmd.Permissions, saved.Permissions)
			assert.Equal(t, cmd.Enabled, saved.Enabled)
		})

		t.Run("successfully updates an existing user", func(t *testing.T) {
			id := uuid.New()
			cmd := newUser()
			cmd.ID = id
			require.NoError(t, repo.Save(t.Context(), cmd))

			cmd.Name = "Updated User"
			cmd.Enabled = false
			err := repo.Save(t.Context(), cmd)
			require.NoError(t, err)

			saved, err := repo.FindByID(t.Context(), id)
			require.NoError(t, err)
			assert.Equal(t, "Updated User", saved.Name)
			assert.False(t, saved.Enabled)
		})
	})

	t.Run("FindByUsername", func(t *testing.T) {
		t.Run("returns user by exact username", func(t *testing.T) {
			cmd := newUser()
			require.NoError(t, repo.Save(t.Context(), cmd))

			saved, err := repo.FindByUsername(t.Context(), cmd.Username)
			require.NoError(t, err)
			require.NotNil(t, saved)
			assert.Equal(t, cmd.ID, saved.ID)
		})

		t.Run("returns nil if not found", func(t *testing.T) {
			saved, err := repo.FindByUsername(t.Context(), "nonexistent")
			require.NoError(t, err)
			assert.Nil(t, saved)
		})
	})

	t.Run("FindPage", func(t *testing.T) {
		t.Run("returns a page of users filtered by name or username", func(t *testing.T) {
			prefix := uuid.New().String()
			names := []string{
				prefix + "Alpha",
				prefix + "Beta",
				prefix + "Gamma",
			}

			for _, name := range names {
				cmd := newUser()
				cmd.ID = uuid.New()
				cmd.Name = name
				cmd.Username = uuid.New().String()
				require.NoError(t, repo.Save(t.Context(), cmd))
			}

			byUsername := newUser()
			byUsername.ID = uuid.New()
			byUsername.Name = "Hidden Name"
			byUsername.Username = prefix + "User"
			require.NoError(t, repo.Save(t.Context(), byUsername))

			other := newUser()
			other.ID = uuid.New()
			other.Name = "Other" + uuid.New().String()
			other.Username = "Other" + uuid.New().String()
			require.NoError(t, repo.Save(t.Context(), other))

			search := prefix
			page, err := repo.FindPage(t.Context(), 10, 0, &search)
			require.NoError(t, err)

			assert.GreaterOrEqual(t, page.TotalItems, 4)

			for _, item := range page.Contents {
				found := false
				if strings.Contains(item.Name, prefix) || strings.Contains(item.Username, prefix) {
					found = true
				}
				assert.True(
					t,
					found,
					"Item %s (user: %s) should match %s",
					item.Name,
					item.Username,
					prefix,
				)
			}
		})
	})

	t.Run("IsEnabledByID", func(t *testing.T) {
		t.Run("returns true when enabled", func(t *testing.T) {
			cmd := newUser()
			cmd.Enabled = true
			require.NoError(t, repo.Save(t.Context(), cmd))

			enabled, err := repo.IsEnabledByID(t.Context(), cmd.ID)
			require.NoError(t, err)
			assert.True(t, enabled)
		})

		t.Run("returns false when disabled", func(t *testing.T) {
			cmd := newUser()
			cmd.Enabled = false
			require.NoError(t, repo.Save(t.Context(), cmd))

			enabled, err := repo.IsEnabledByID(t.Context(), cmd.ID)
			require.NoError(t, err)
			assert.False(t, enabled)
		})

		t.Run("returns false when not exists", func(t *testing.T) {
			enabled, err := repo.IsEnabledByID(t.Context(), uuid.New())
			require.NoError(t, err)
			assert.False(t, enabled)
		})
	})

	t.Run("Count", func(t *testing.T) {
		t.Run("returns total user count", func(t *testing.T) {
			initial, err := repo.Count(t.Context())
			require.NoError(t, err)

			cmd := newUser()
			require.NoError(t, repo.Save(t.Context(), cmd))

			cmd2 := newUser()
			cmd2.ID = uuid.New()
			cmd2.Username = uuid.New().String()
			require.NoError(t, repo.Save(t.Context(), cmd2))

			current, err := repo.Count(t.Context())
			require.NoError(t, err)
			assert.Equal(t, initial+2, current)
		})
	})

	t.Run("DeleteByID", func(t *testing.T) {
		t.Run("removes the user", func(t *testing.T) {
			cmd := newUser()
			require.NoError(t, repo.Save(t.Context(), cmd))

			err := repo.DeleteByID(t.Context(), cmd.ID)
			require.NoError(t, err)

			saved, err := repo.FindByID(t.Context(), cmd.ID)
			require.NoError(t, err)
			assert.Nil(t, saved)
		})
	})
}
