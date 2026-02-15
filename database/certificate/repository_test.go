package certificate

import (
	"testing"
	"time"

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
		t.Run("successfully saves a new certificate", func(t *testing.T) {
			cmd := newCertificate()

			err := repo.Save(t.Context(), cmd)
			require.NoError(t, err)

			saved, err := repo.FindByID(t.Context(), cmd.ID)
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
			require.NoError(t, repo.Save(t.Context(), cmd))

			cmd.Metadata = new("Updated Metadata")
			cmd.RenewAfter = new(time.Now().Add(48 * time.Hour))
			err := repo.Save(t.Context(), cmd)
			require.NoError(t, err)

			saved, err := repo.FindByID(t.Context(), id)
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
				require.NoError(t, repo.Save(t.Context(), cmd))
			}

			other := newCertificate()
			other.ID = uuid.New()
			other.DomainNames = []string{"other.com"}
			require.NoError(t, repo.Save(t.Context(), other))

			page, err := repo.FindPage(t.Context(), 0, 10, new(prefix))
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
			require.NoError(t, repo.Save(t.Context(), cmd))

			err := repo.DeleteByID(t.Context(), cmd.ID)
			require.NoError(t, err)

			exists, err := repo.ExistsByID(t.Context(), cmd.ID)
			require.NoError(t, err)
			assert.False(t, exists)
		})
	})

	t.Run("FindAllDueToRenew", func(t *testing.T) {
		t.Run("returns certificates eligible for rotation", func(t *testing.T) {
			cmd := newCertificate()
			cmd.RenewAfter = new(time.Now().Add(-1 * time.Hour))
			require.NoError(t, repo.Save(t.Context(), cmd))

			candidates, err := repo.FindAllDueToRenew(t.Context())
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
			cmd.RenewAfter = new(time.Now().Add(1 * time.Hour))
			require.NoError(t, repo.Save(t.Context(), cmd))

			candidates, err := repo.FindAllDueToRenew(t.Context())
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
			require.NoError(t, repo.Save(t.Context(), cmd))

			inUse, err := repo.InUseByID(t.Context(), cmd.ID)
			require.NoError(t, err)
			assert.False(t, inUse)
		})
	})

	t.Run("GetAutoRenewSettings", func(t *testing.T) {
		t.Run("returns settings", func(t *testing.T) {
			settings, err := repo.GetAutoRenewSettings(t.Context())
			require.NoError(t, err)
			require.NotNil(t, settings)
		})
	})
}
