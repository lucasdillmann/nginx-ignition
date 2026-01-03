package accesslist

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"

	"dillmann.com.br/nginx-ignition/core/accesslist"
	"dillmann.com.br/nginx-ignition/database/common/database"
	"dillmann.com.br/nginx-ignition/database/common/testutils"
)

type hostModel struct {
	bun.BaseModel `bun:"host"`

	DomainNames         []string  `bun:"domain_names,array"`
	ID                  uuid.UUID `bun:"id,pk"`
	AccessListID        uuid.UUID `bun:"access_list_id"`
	Enabled             bool      `bun:"enabled,notnull"`
	DefaultServer       bool      `bun:"default_server,notnull"`
	WebsocketSupport    bool      `bun:"websocket_support,notnull"`
	HTTP2Support        bool      `bun:"http2_support,notnull"`
	RedirectHTTPToHTTPS bool      `bun:"redirect_http_to_https,notnull"`
	UseGlobalBindings   bool      `bun:"use_global_bindings,notnull"`
}

func Test_Repository(t *testing.T) {
	testutils.RunWithMockedDatabases(t, runRepositoryTests)
}

func runRepositoryTests(t *testing.T, db *database.Database) {
	ctx := context.Background()
	repo := New(db)

	t.Run("Save", func(t *testing.T) {
		t.Run("successfully saves a new access list", func(t *testing.T) {
			cmd := newAccessList()

			err := repo.Save(ctx, cmd)
			require.NoError(t, err)

			exists, err := repo.ExistsByID(ctx, cmd.ID)
			require.NoError(t, err)
			assert.True(t, exists)
		})

		t.Run("successfully updates an existing access list", func(t *testing.T) {
			id := uuid.New()
			cmd := newAccessList()
			cmd.ID = id
			require.NoError(t, repo.Save(ctx, cmd))

			cmd.Name = "Updated Name"
			err := repo.Save(ctx, cmd)
			require.NoError(t, err)

			found, err := repo.FindByID(ctx, id)
			require.NoError(t, err)
			require.NotNil(t, found)
			assert.Equal(t, "Updated Name", found.Name)
		})
	})

	t.Run("FindByID", func(t *testing.T) {
		t.Run("returns the access list when it exists", func(t *testing.T) {
			cmd := newAccessList()
			require.NoError(t, repo.Save(ctx, cmd))

			found, err := repo.FindByID(ctx, cmd.ID)
			require.NoError(t, err)
			require.NotNil(t, found)
			assert.Equal(t, cmd.Name, found.Name)
			assert.Equal(t, cmd.Realm, found.Realm)
			assert.Len(t, found.Entries, len(cmd.Entries))
			assert.Len(t, found.Credentials, len(cmd.Credentials))
		})

		t.Run("returns nil when the access list does not exist", func(t *testing.T) {
			found, err := repo.FindByID(ctx, uuid.New())
			assert.NoError(t, err)
			assert.Nil(t, found)
		})
	})

	t.Run("ExistsByID", func(t *testing.T) {
		t.Run("returns true when it exists", func(t *testing.T) {
			cmd := newAccessList()
			require.NoError(t, repo.Save(ctx, cmd))

			exists, err := repo.ExistsByID(ctx, cmd.ID)
			assert.NoError(t, err)
			assert.True(t, exists)
		})

		t.Run("returns false when it does not exist", func(t *testing.T) {
			exists, err := repo.ExistsByID(ctx, uuid.New())
			assert.NoError(t, err)
			assert.False(t, exists)
		})
	})

	t.Run("DeleteByID", func(t *testing.T) {
		t.Run("successfully deletes an existing access list", func(t *testing.T) {
			cmd := newAccessList()
			require.NoError(t, repo.Save(ctx, cmd))

			err := repo.DeleteByID(ctx, cmd.ID)
			assert.NoError(t, err)

			exists, _ := repo.ExistsByID(ctx, cmd.ID)
			assert.False(t, exists)
		})
	})

	t.Run("InUseByID", func(t *testing.T) {
		t.Run("returns true when used by a host", func(t *testing.T) {
			cmd := newAccessList()
			require.NoError(t, repo.Save(ctx, cmd))

			host := &hostModel{
				ID:           uuid.New(),
				Enabled:      true,
				DomainNames:  []string{},
				AccessListID: cmd.ID,
			}

			_, err := db.Insert().Model(host).Exec(ctx)
			require.NoError(t, err)

			inUse, err := repo.InUseByID(ctx, cmd.ID)
			assert.NoError(t, err)
			assert.True(t, inUse)
		})

		t.Run("returns false when not in use", func(t *testing.T) {
			inUse, err := repo.InUseByID(ctx, uuid.New())
			assert.NoError(t, err)
			assert.False(t, inUse)
		})
	})

	t.Run("FindPage", func(t *testing.T) {
		t.Run("returns a page of access lists", func(t *testing.T) {
			for i := 0; i < 3; i++ {
				cmd := newAccessList()
				cmd.ID = uuid.New()
				require.NoError(t, repo.Save(ctx, cmd))
			}

			page, err := repo.FindPage(ctx, 0, 2, nil)
			require.NoError(t, err)
			assert.Equal(t, 2, len(page.Contents))
			assert.GreaterOrEqual(t, page.TotalItems, 3)
		})

		t.Run("filters by search terms", func(t *testing.T) {
			cmd := newAccessList()
			cmd.Name = "SearchMe"
			require.NoError(t, repo.Save(ctx, cmd))

			search := "SearchMe"
			page, err := repo.FindPage(ctx, 0, 10, &search)
			require.NoError(t, err)
			assert.GreaterOrEqual(t, page.TotalItems, 1)
		})
	})

	t.Run("FindAll", func(t *testing.T) {
		t.Run("returns all access lists", func(t *testing.T) {
			cmd := newAccessList()
			require.NoError(t, repo.Save(ctx, cmd))

			all, err := repo.FindAll(ctx)
			assert.NoError(t, err)
			assert.GreaterOrEqual(t, len(all), 1)
		})
	})
}

func newAccessList() *accesslist.AccessList {
	return &accesslist.AccessList{
		ID:             uuid.New(),
		Name:           uuid.NewString(),
		Realm:          "Restricted Area",
		SatisfyAll:     false,
		DefaultOutcome: accesslist.AllowOutcome,
		Entries: []accesslist.Entry{
			{
				Outcome:       accesslist.DenyOutcome,
				SourceAddress: []string{"192.168.1.1"},
				Priority:      1,
			},
		},
		Credentials: []accesslist.Credentials{
			{
				Username: "user",
				Password: "password",
			},
		},
	}
}
