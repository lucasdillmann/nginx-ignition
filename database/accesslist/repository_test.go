package accesslist

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
		t.Run("successfully saves a new access list", func(t *testing.T) {
			cmd := newAccessList()

			err := repo.Save(t.Context(), cmd)
			require.NoError(t, err)

			exists, err := repo.ExistsByID(t.Context(), cmd.ID)
			require.NoError(t, err)
			assert.True(t, exists)
		})

		t.Run("successfully updates an existing access list", func(t *testing.T) {
			id := uuid.New()
			cmd := newAccessList()
			cmd.ID = id
			require.NoError(t, repo.Save(t.Context(), cmd))

			cmd.Name = "Updated Name"
			err := repo.Save(t.Context(), cmd)
			require.NoError(t, err)

			found, err := repo.FindByID(t.Context(), id)
			require.NoError(t, err)
			require.NotNil(t, found)
			assert.Equal(t, "Updated Name", found.Name)
		})
	})

	t.Run("FindByID", func(t *testing.T) {
		t.Run("returns the access list when it exists", func(t *testing.T) {
			cmd := newAccessList()
			require.NoError(t, repo.Save(t.Context(), cmd))

			found, err := repo.FindByID(t.Context(), cmd.ID)
			require.NoError(t, err)
			require.NotNil(t, found)
			assert.Equal(t, cmd.Name, found.Name)
			assert.Equal(t, cmd.Realm, found.Realm)
			assert.Len(t, found.Entries, len(cmd.Entries))
			assert.Len(t, found.Credentials, len(cmd.Credentials))
		})

		t.Run("returns nil when the access list does not exist", func(t *testing.T) {
			found, err := repo.FindByID(t.Context(), uuid.New())
			assert.NoError(t, err)
			assert.Nil(t, found)
		})
	})

	t.Run("ExistsByID", func(t *testing.T) {
		t.Run("returns true when it exists", func(t *testing.T) {
			cmd := newAccessList()
			require.NoError(t, repo.Save(t.Context(), cmd))

			exists, err := repo.ExistsByID(t.Context(), cmd.ID)
			assert.NoError(t, err)
			assert.True(t, exists)
		})

		t.Run("returns false when it does not exist", func(t *testing.T) {
			exists, err := repo.ExistsByID(t.Context(), uuid.New())
			assert.NoError(t, err)
			assert.False(t, exists)
		})
	})

	t.Run("DeleteByID", func(t *testing.T) {
		t.Run("successfully deletes an existing access list", func(t *testing.T) {
			cmd := newAccessList()
			require.NoError(t, repo.Save(t.Context(), cmd))

			err := repo.DeleteByID(t.Context(), cmd.ID)
			assert.NoError(t, err)

			exists, _ := repo.ExistsByID(t.Context(), cmd.ID)
			assert.False(t, exists)
		})
	})

	t.Run("InUseByID", func(t *testing.T) {
		t.Run("returns true when used by a host", func(t *testing.T) {
			cmd := newAccessList()
			require.NoError(t, repo.Save(t.Context(), cmd))

			host := &hostModel{
				ID:           uuid.New(),
				Enabled:      true,
				StatsEnabled: false,
				DomainNames:  []string{},
				AccessListID: cmd.ID,
			}

			_, err := db.Insert().Model(host).Exec(t.Context())
			require.NoError(t, err)

			inUse, err := repo.InUseByID(t.Context(), cmd.ID)
			assert.NoError(t, err)
			assert.True(t, inUse)
		})

		t.Run("returns false when not in use", func(t *testing.T) {
			inUse, err := repo.InUseByID(t.Context(), uuid.New())
			assert.NoError(t, err)
			assert.False(t, inUse)
		})
	})

	t.Run("FindPage", func(t *testing.T) {
		t.Run("returns a page of access lists", func(t *testing.T) {
			for i := 0; i < 3; i++ {
				cmd := newAccessList()
				cmd.ID = uuid.New()
				require.NoError(t, repo.Save(t.Context(), cmd))
			}

			page, err := repo.FindPage(t.Context(), 0, 2, nil)
			require.NoError(t, err)
			assert.Equal(t, 2, len(page.Contents))
			assert.GreaterOrEqual(t, page.TotalItems, 3)
		})

		t.Run("filters by search terms", func(t *testing.T) {
			cmd := newAccessList()
			cmd.Name = "SearchMe"
			require.NoError(t, repo.Save(t.Context(), cmd))

			search := "SearchMe"
			page, err := repo.FindPage(t.Context(), 0, 10, &search)
			require.NoError(t, err)
			assert.GreaterOrEqual(t, page.TotalItems, 1)
		})
	})

	t.Run("FindAll", func(t *testing.T) {
		t.Run("returns all access lists", func(t *testing.T) {
			cmd := newAccessList()
			require.NoError(t, repo.Save(t.Context(), cmd))

			all, err := repo.FindAll(t.Context())
			assert.NoError(t, err)
			assert.GreaterOrEqual(t, len(all), 1)
		})
	})
}
