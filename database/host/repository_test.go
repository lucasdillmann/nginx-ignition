package host

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"dillmann.com.br/nginx-ignition/core/binding"
	"dillmann.com.br/nginx-ignition/core/common/ptr"
	"dillmann.com.br/nginx-ignition/core/host"
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
		t.Run("successfully saves a new host", func(t *testing.T) {
			cmd := newHost()

			err := repo.Save(ctx, cmd)
			require.NoError(t, err)

			saved, err := repo.FindByID(ctx, cmd.ID)
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
			require.NoError(t, repo.Save(ctx, cmd))

			cmd.Enabled = false
			cmd.DomainNames = []string{"updated.example.com"}
			err := repo.Save(ctx, cmd)
			require.NoError(t, err)

			saved, err := repo.FindByID(ctx, id)
			require.NoError(t, err)
			assert.ElementsMatch(t, []string{"updated.example.com"}, saved.DomainNames)
			assert.False(t, saved.Enabled)
		})
	})

	t.Run("ExistsByID", func(t *testing.T) {
		t.Run("returns true when exists", func(t *testing.T) {
			cmd := newHost()
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
				require.NoError(t, repo.Save(ctx, cmd))
			}

			other := newHost()
			other.ID = uuid.New()
			other.DomainNames = []string{"other.com"}
			require.NoError(t, repo.Save(ctx, other))

			search := prefix
			page, err := repo.FindPage(ctx, 10, 0, &search)
			require.NoError(t, err)

			assert.GreaterOrEqual(t, page.TotalItems, 3)

			for _, item := range page.Contents {
				found := false
				for _, d := range item.DomainNames {
					if contains(d, prefix) {
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
			require.NoError(t, repo.Save(ctx, enabled))

			disabled := newHost()
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

	t.Run("FindDefault", func(t *testing.T) {
		t.Run("returns the default server", func(t *testing.T) {
			cleanup(ctx, t, repo)

			def := newHost()
			def.ID = uuid.New()
			def.DefaultServer = true
			def.DomainNames = nil
			require.NoError(t, repo.Save(ctx, def))

			normal := newHost()
			normal.ID = uuid.New()
			normal.DefaultServer = false
			require.NoError(t, repo.Save(ctx, normal))

			found, err := repo.FindDefault(ctx)
			require.NoError(t, err)
			require.NotNil(t, found)
			assert.Equal(t, def.ID, found.ID)
		})

		t.Run("returns nil when no default server", func(t *testing.T) {
			cleanup(ctx, t, repo)

			existing, err := repo.FindDefault(ctx)
			require.NoError(t, err)
			if existing != nil {
				err = repo.DeleteByID(ctx, existing.ID)
				require.NoError(t, err)
			}

			found, err := repo.FindDefault(ctx)
			require.NoError(t, err)
			assert.Nil(t, found)
		})
	})

	t.Run("DeleteByID", func(t *testing.T) {
		t.Run("removes the host", func(t *testing.T) {
			cmd := newHost()
			require.NoError(t, repo.Save(ctx, cmd))

			err := repo.DeleteByID(ctx, cmd.ID)
			require.NoError(t, err)

			exists, err := repo.ExistsByID(ctx, cmd.ID)
			require.NoError(t, err)
			assert.False(t, exists)
		})
	})
}

func cleanup(ctx context.Context, t *testing.T, repo host.Repository) {
	result, err := repo.(*repository).database.Unwrap().Query("SELECT id FROM host")
	require.NoError(t, err)

	defer result.Close()

	ids := make([]string, 0)
	for result.Next() {
		var id string

		err = result.Scan(&id)
		require.NoError(t, err)

		ids = append(ids, id)
	}

	result.Close()

	for _, id := range ids {
		err = repo.DeleteByID(ctx, uuid.MustParse(id))
		require.NoError(t, err)
	}
}

func newHost() *host.Host {
	return &host.Host{
		ID:                uuid.New(),
		DomainNames:       []string{"example.com"},
		Enabled:           true,
		DefaultServer:     false,
		UseGlobalBindings: true,
		Routes: []host.Route{
			{
				ID:         uuid.New(),
				Priority:   10,
				Type:       host.StaticResponseRouteType,
				SourcePath: "/",
				Settings: host.RouteSettings{
					IncludeForwardHeaders:   true,
					ProxySSLServerName:      false,
					KeepOriginalDomainName:  true,
					DirectoryListingEnabled: false,
					Custom:                  ptr.Of("# Custom config"),
				},
				Response: &host.RouteStaticResponse{
					StatusCode: 200,
					Payload:    ptr.Of("OK"),
				},
			},
		},
		Bindings: []binding.Binding{
			{
				ID:   uuid.New(),
				Type: binding.HTTPBindingType,
				IP:   "0.0.0.0",
				Port: 8080,
			},
		},
		FeatureSet: host.FeatureSet{
			WebsocketSupport:    true,
			HTTP2Support:        true,
			RedirectHTTPToHTTPS: false,
		},
		VPNs: []host.VPN{},
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[0:len(substr)] == substr
}
