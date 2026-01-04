package settings

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"dillmann.com.br/nginx-ignition/core/binding"
	"dillmann.com.br/nginx-ignition/core/common/ptr"
	"dillmann.com.br/nginx-ignition/core/settings"
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
		t.Run("successfully saves settings", func(t *testing.T) {
			cmd := newSettings()

			err := repo.Save(ctx, cmd)
			require.NoError(t, err)

			saved, err := repo.Get(ctx)
			require.NoError(t, err)
			require.NotNil(t, saved)

			assert.Equal(t, cmd.Nginx, saved.Nginx)
			assert.Equal(t, cmd.LogRotation, saved.LogRotation)
			assert.Equal(t, cmd.CertificateAutoRenew, saved.CertificateAutoRenew)
			assert.ElementsMatch(t, cmd.GlobalBindings, saved.GlobalBindings)
		})

		t.Run("successfully updates existing settings", func(t *testing.T) {
			cmd := newSettings()
			require.NoError(t, repo.Save(ctx, cmd))

			cmd.Nginx.WorkerProcesses = 4
			cmd.Nginx.GzipEnabled = false
			err := repo.Save(ctx, cmd)
			require.NoError(t, err)

			saved, err := repo.Get(ctx)
			require.NoError(t, err)
			assert.Equal(t, 4, saved.Nginx.WorkerProcesses)
			assert.False(t, saved.Nginx.GzipEnabled)
		})
	})
}

func newSettings() *settings.Settings {
	return &settings.Settings{
		Nginx: &settings.NginxSettings{
			Timeouts: &settings.NginxTimeoutsSettings{
				Read:       60,
				Connect:    60,
				Send:       60,
				Keepalive:  75,
				ClientBody: 60,
			},
			Buffers: &settings.NginxBuffersSettings{
				LargeClientHeader: &settings.NginxBufferSize{
					SizeKb: 8,
					Amount: 4,
				},
				Output: &settings.NginxBufferSize{
					SizeKb: 32,
					Amount: 4,
				},
				ClientBodyKb:   16,
				ClientHeaderKb: 1,
			},
			Logs: &settings.NginxLogsSettings{
				ServerLogsLevel:   settings.WarnLogLevel,
				ErrorLogsLevel:    settings.ErrorLogLevel,
				ServerLogsEnabled: true,
				AccessLogsEnabled: true,
				ErrorLogsEnabled:  true,
			},
			Custom:              ptr.Of("# Custom Nginx Config"),
			RuntimeUser:         "nginx",
			DefaultContentType:  "text/html",
			WorkerProcesses:     1,
			WorkerConnections:   1024,
			MaximumBodySizeMb:   10,
			ServerTokensEnabled: false,
			TCPNoDelayEnabled:   true,
			GzipEnabled:         true,
			SendfileEnabled:     true,
		},
		LogRotation: &settings.LogRotationSettings{
			IntervalUnit:      settings.DaysTimeUnit,
			MaximumLines:      10000,
			IntervalUnitCount: 7,
			Enabled:           true,
		},
		CertificateAutoRenew: &settings.CertificateAutoRenewSettings{
			IntervalUnit:      settings.DaysTimeUnit,
			IntervalUnitCount: 30,
			Enabled:           true,
		},
		GlobalBindings: []binding.Binding{
			{
				Type: binding.HTTPBindingType,
				IP:   "0.0.0.0",
				Port: 80,
			},
		},
	}
}
