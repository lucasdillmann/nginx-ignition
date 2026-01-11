package settings

import (
	"testing"

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
		t.Run("successfully saves settings", func(t *testing.T) {
			cmd := newSettings()

			err := repo.Save(t.Context(), cmd)
			require.NoError(t, err)

			saved, err := repo.Get(t.Context())
			require.NoError(t, err)
			require.NotNil(t, saved)

			assert.Equal(t, cmd.Nginx, saved.Nginx)
			assert.Equal(t, cmd.LogRotation, saved.LogRotation)
			assert.Equal(t, cmd.CertificateAutoRenew, saved.CertificateAutoRenew)
			assert.ElementsMatch(t, cmd.GlobalBindings, saved.GlobalBindings)
		})

		t.Run("successfully updates existing settings", func(t *testing.T) {
			cmd := newSettings()
			require.NoError(t, repo.Save(t.Context(), cmd))

			cmd.Nginx.WorkerProcesses = 4
			cmd.Nginx.GzipEnabled = false
			err := repo.Save(t.Context(), cmd)
			require.NoError(t, err)

			saved, err := repo.Get(t.Context())
			require.NoError(t, err)
			assert.Equal(t, 4, saved.Nginx.WorkerProcesses)
			assert.False(t, saved.Nginx.GzipEnabled)
		})
	})
}
