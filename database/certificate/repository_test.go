package certificate

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"dillmann.com.br/nginx-ignition/core/common/ptr"
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
		t.Run("successfully saves a new certificate", func(t *testing.T) {
			cmd := newCertificate()

			err := repo.Save(ctx, cmd)
			require.NoError(t, err)

			saved, err := repo.FindByID(ctx, cmd.ID)
			require.NoError(t, err)
			require.NotNil(t, saved)
			assert.Equal(t, cmd.ProviderID, saved.ProviderID)
			assert.Equal(t, cmd.PrivateKey, saved.PrivateKey)
			assert.Equal(t, cmd.PublicKey, saved.PublicKey)
			assert.Equal(t, cmd.Parameters, saved.Parameters)
			assert.Equal(t, cmd.Metadata, saved.Metadata)
			assert.WithinDuration(t, cmd.IssuedAt, saved.IssuedAt, time.Second)
			assert.WithinDuration(t, cmd.ValidFrom, saved.ValidFrom, time.Second)
			assert.WithinDuration(t, cmd.ValidUntil, saved.ValidUntil, time.Second)
			assert.WithinDuration(t, *cmd.RenewAfter, *saved.RenewAfter, time.Second)
			assert.ElementsMatch(t, cmd.DomainNames, saved.DomainNames)
			assert.ElementsMatch(t, cmd.CertificationChain, saved.CertificationChain)
		})

		t.Run("successfully updates an existing certificate", func(t *testing.T) {
			id := uuid.New()
			cmd := newCertificate()
			cmd.ID = id
			require.NoError(t, repo.Save(ctx, cmd))

			cmd.Metadata = ptr.Of("Updated Metadata")
			cmd.RenewAfter = ptr.Of(time.Now().Add(48 * time.Hour))
			err := repo.Save(ctx, cmd)
			require.NoError(t, err)

			saved, err := repo.FindByID(ctx, id)
			require.NoError(t, err)
			assert.Equal(t, "Updated Metadata", *saved.Metadata)
			assert.WithinDuration(t, *cmd.RenewAfter, *saved.RenewAfter, time.Second)
		})
	})

	t.Run("FindPage", func(t *testing.T) {
		t.Run("returns a page of certificates filtered by domain name", func(t *testing.T) {
			prefix := uuid.New().String()
			domains := []string{
				prefix + ".com",
				prefix + ".org",
				prefix + ".net",
			}

			for _, domain := range domains {
				cmd := newCertificate()
				cmd.ID = uuid.New()
				cmd.DomainNames = []string{domain}
				require.NoError(t, repo.Save(ctx, cmd))
			}

			other := newCertificate()
			other.ID = uuid.New()
			other.DomainNames = []string{"other.com"}
			require.NoError(t, repo.Save(ctx, other))

			search := prefix
			page, err := repo.FindPage(ctx, 0, 10, &search)
			require.NoError(t, err)

			assert.GreaterOrEqual(t, page.TotalItems, 3)

			for _, item := range page.Contents {
				found := false
				for _, d := range item.DomainNames {
					if d == domains[0] || d == domains[1] || d == domains[2] {
						found = true
						break
					}
				}
				assert.True(t, found)
			}
		})
	})

	t.Run("DeleteByID", func(t *testing.T) {
		t.Run("removes the certificate", func(t *testing.T) {
			cmd := newCertificate()
			require.NoError(t, repo.Save(ctx, cmd))

			err := repo.DeleteByID(ctx, cmd.ID)
			require.NoError(t, err)

			exists, err := repo.ExistsByID(ctx, cmd.ID)
			require.NoError(t, err)
			assert.False(t, exists)
		})
	})

	t.Run("FindAllDueToRenew", func(t *testing.T) {
		t.Run("returns certificates eligible for rotation", func(t *testing.T) {
			cmd := newCertificate()
			cmd.RenewAfter = ptr.Of(time.Now().Add(-1 * time.Hour))
			require.NoError(t, repo.Save(ctx, cmd))

			candidates, err := repo.FindAllDueToRenew(ctx)
			require.NoError(t, err)

			found := false
			for _, c := range candidates {
				if c.ID == cmd.ID {
					found = true
					break
				}
			}
			assert.True(t, found)
		})

		t.Run("does not return certificates not yet eligible", func(t *testing.T) {
			cmd := newCertificate()
			cmd.RenewAfter = ptr.Of(time.Now().Add(1 * time.Hour))
			require.NoError(t, repo.Save(ctx, cmd))

			candidates, err := repo.FindAllDueToRenew(ctx)
			require.NoError(t, err)

			found := false
			for _, c := range candidates {
				if c.ID == cmd.ID {
					found = true
					break
				}
			}
			assert.False(t, found)
		})
	})

	t.Run("InUseByID", func(t *testing.T) {
		t.Run("returns false when not in use", func(t *testing.T) {
			cmd := newCertificate()
			require.NoError(t, repo.Save(ctx, cmd))

			inUse, err := repo.InUseByID(ctx, cmd.ID)
			require.NoError(t, err)
			assert.False(t, inUse)
		})
	})

	t.Run("GetAutoRenewSettings", func(t *testing.T) {
		t.Run("returns settings", func(t *testing.T) {
			settings, err := repo.GetAutoRenewSettings(ctx)
			require.NoError(t, err)
			require.NotNil(t, settings)
		})
	})
}
