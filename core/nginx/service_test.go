package nginx

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"dillmann.com.br/nginx-ignition/core/common/configuration"
	"dillmann.com.br/nginx-ignition/core/common/coreerror"
	"dillmann.com.br/nginx-ignition/core/host"
	"dillmann.com.br/nginx-ignition/core/settings"
)

func Test_service(t *testing.T) {
	t.Run("GetMainLogs", func(t *testing.T) {
		ctx := context.Background()
		tmpDir := t.TempDir()
		logsDir := filepath.Join(tmpDir, "logs")
		err := os.Mkdir(logsDir, 0o755)
		require.NoError(t, err)

		mainLogPath := filepath.Join(logsDir, "main.log")
		err = os.WriteFile(mainLogPath, []byte("line1\nline2\nline3\n"), 0o644)
		require.NoError(t, err)

		cfg := configuration.NewWithOverrides(map[string]string{
			"nginx-ignition.nginx.config-path": tmpDir,
		})
		nginxService := &service{
			logReader: newLogReader(cfg),
		}

		t.Run("returns requested number of lines in reverse order", func(t *testing.T) {
			lines, err := nginxService.GetMainLogs(ctx, 2)
			assert.NoError(t, err)
			assert.Equal(t, []string{"line3", "line2"}, lines)
		})
	})

	t.Run("GetHostLogs", func(t *testing.T) {
		ctx := context.Background()
		tmpDir := t.TempDir()
		logsDir := filepath.Join(tmpDir, "logs")
		err := os.Mkdir(logsDir, 0o755)
		require.NoError(t, err)

		hostID := uuid.New()
		hostLogPath := filepath.Join(logsDir, fmt.Sprintf("host-%s.access.log", hostID))
		err = os.WriteFile(hostLogPath, []byte("access1\naccess2\n"), 0o644)
		require.NoError(t, err)

		cfg := configuration.NewWithOverrides(map[string]string{
			"nginx-ignition.nginx.config-path": tmpDir,
		})
		nginxService := &service{
			logReader: newLogReader(cfg),
		}

		t.Run("returns host specific logs", func(t *testing.T) {
			lines, err := nginxService.GetHostLogs(ctx, hostID, "access", 1)
			assert.NoError(t, err)
			assert.Equal(t, []string{"access2"}, lines)
		})
	})

	t.Run("rotateLogs", func(t *testing.T) {
		ctx := context.Background()
		tmpDir := t.TempDir()
		logsDir := filepath.Join(tmpDir, "logs")
		err := os.Mkdir(logsDir, 0o755)
		require.NoError(t, err)

		fakeNginx := filepath.Join(tmpDir, "nginx_fake")
		err = os.WriteFile(fakeNginx, []byte("#!/bin/sh\nexit 0"), 0o755)
		require.NoError(t, err)

		cfg := configuration.NewWithOverrides(map[string]string{
			"nginx-ignition.nginx.config-path": tmpDir,
		})

		mainLogPath := filepath.Join(logsDir, "main.log")
		err = os.WriteFile(mainLogPath, []byte("line1\nline2\nline3\n"), 0o644)
		require.NoError(t, err)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		settingsCmds := settings.NewMockedCommands(ctrl)
		settingsCmds.EXPECT().Get(ctx).Return(&settings.Settings{
			LogRotation: &settings.LogRotationSettings{
				Enabled:      true,
				MaximumLines: 2,
			},
		}, nil)

		hostCmds := host.NewMockedCommands(ctrl)
		hostCmds.EXPECT().GetAllEnabled(ctx).Return([]host.Host{}, nil)

		nginxService := &service{
			logRotator: newLogRotator(
				cfg,
				settingsCmds,
				hostCmds,
				&processManager{
					binaryPath: fakeNginx,
					configPath: tmpDir,
				},
			),
		}

		err = nginxService.rotateLogs(ctx)
		assert.NoError(t, err)

		content, err := os.ReadFile(mainLogPath)
		require.NoError(t, err)
		assert.Equal(t, "line2\nline3\n", string(content))
	})

	t.Run("Reload", func(t *testing.T) {
		ctx := context.Background()

		t.Run("returns error when not running and failIfNotRunning is true", func(t *testing.T) {
			nginxService := &service{
				semaphore: &semaphore{
					state: stoppedState,
				},
			}

			err := nginxService.Reload(ctx, true)
			assert.Error(t, err)
			var coreErr *coreerror.CoreError
			assert.ErrorAs(t, err, &coreErr)
			assert.Contains(t, coreErr.Message, "not running")
		})
	})

	t.Run("Start", func(t *testing.T) {
		ctx := context.Background()

		t.Run("returns nil if already running", func(t *testing.T) {
			nginxService := &service{
				semaphore: &semaphore{
					state: runningState,
				},
			}

			err := nginxService.Start(ctx)
			assert.NoError(t, err)
		})
	})

	t.Run("Stop", func(t *testing.T) {
		ctx := context.Background()

		t.Run("returns nil if already stopped", func(t *testing.T) {
			nginxService := &service{
				semaphore: &semaphore{
					state: stoppedState,
				},
			}

			err := nginxService.Stop(ctx)
			assert.NoError(t, err)
		})
	})

	t.Run("GetStatus", func(t *testing.T) {
		t.Run("returns true when running", func(t *testing.T) {
			nginxService := &service{
				semaphore: &semaphore{
					state: runningState,
				},
			}
			assert.True(t, nginxService.GetStatus(context.Background()))
		})

		t.Run("returns false when stopped", func(t *testing.T) {
			nginxService := &service{
				semaphore: &semaphore{
					state: stoppedState,
				},
			}
			assert.False(t, nginxService.GetStatus(context.Background()))
		})
	})
}
