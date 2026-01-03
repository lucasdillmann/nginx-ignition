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

func Test_Service_GetMainLogs(t *testing.T) {
	ctx := context.Background()
	tmpDir := t.TempDir()
	logsDir := filepath.Join(tmpDir, "logs")
	err := os.Mkdir(logsDir, 0o755)
	require.NoError(t, err)

	mainLogPath := filepath.Join(logsDir, "main.log")
	err = os.WriteFile(mainLogPath, []byte("line1\nline2\nline3\n"), 0o644)
	require.NoError(t, err)

	t.Setenv("NGINX_IGNITION_NGINX_CONFIG_PATH", tmpDir)
	cfg := configuration.New()
	s := &service{
		logReader: newLogReader(cfg),
	}

	t.Run("returns requested number of lines in reverse order", func(t *testing.T) {
		lines, err := s.GetMainLogs(ctx, 2)
		assert.NoError(t, err)
		assert.Equal(t, []string{"line3", "line2"}, lines)
	})
}

func Test_Service_GetHostLogs(t *testing.T) {
	ctx := context.Background()
	tmpDir := t.TempDir()
	logsDir := filepath.Join(tmpDir, "logs")
	err := os.Mkdir(logsDir, 0o755)
	require.NoError(t, err)

	hostID := uuid.New()
	hostLogPath := filepath.Join(logsDir, fmt.Sprintf("host-%s.access.log", hostID))
	err = os.WriteFile(hostLogPath, []byte("access1\naccess2\n"), 0o644)
	require.NoError(t, err)

	t.Setenv("NGINX_IGNITION_NGINX_CONFIG_PATH", tmpDir)
	cfg := configuration.New()
	s := &service{
		logReader: newLogReader(cfg),
	}

	t.Run("returns host specific logs", func(t *testing.T) {
		lines, err := s.GetHostLogs(ctx, hostID, "access", 1)
		assert.NoError(t, err)
		assert.Equal(t, []string{"access2"}, lines)
	})
}

func Test_Service_RotateLogs(t *testing.T) {
	ctx := context.Background()
	tmpDir := t.TempDir()
	logsDir := filepath.Join(tmpDir, "logs")
	err := os.Mkdir(logsDir, 0o755)
	require.NoError(t, err)

	fakeNginx := filepath.Join(tmpDir, "nginx_fake")
	err = os.WriteFile(fakeNginx, []byte("#!/bin/sh\nexit 0"), 0o755)
	require.NoError(t, err)

	t.Setenv("NGINX_IGNITION_NGINX_CONFIG_PATH", tmpDir)
	cfg := configuration.New()

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

	s := &service{
		logRotator: newLogRotator(
			cfg,
			settingsCmds,
			hostCmds,
			&processManager{binaryPath: fakeNginx, configPath: tmpDir},
		),
	}

	err = s.rotateLogs(ctx)
	assert.NoError(t, err)

	content, err := os.ReadFile(mainLogPath)
	require.NoError(t, err)
	assert.Equal(t, "line2\nline3\n", string(content))
}

func Test_Service_Reload(t *testing.T) {
	ctx := context.Background()

	t.Run("returns error when not running and failIfNotRunning is true", func(t *testing.T) {
		s := &service{
			semaphore: &semaphore{
				state: stoppedState,
			},
		}

		err := s.Reload(ctx, true)
		assert.Error(t, err)
		var coreErr *coreerror.CoreError
		assert.ErrorAs(t, err, &coreErr)
		assert.Contains(t, coreErr.Message, "not running")
	})
}

func Test_Service_Start(t *testing.T) {
	ctx := context.Background()

	t.Run("returns nil if already running", func(t *testing.T) {
		s := &service{
			semaphore: &semaphore{
				state: runningState,
			},
		}

		err := s.Start(ctx)
		assert.NoError(t, err)
	})
}

func Test_Service_Stop(t *testing.T) {
	ctx := context.Background()

	t.Run("returns nil if already stopped", func(t *testing.T) {
		s := &service{
			semaphore: &semaphore{
				state: stoppedState,
			},
		}

		err := s.Stop(ctx)
		assert.NoError(t, err)
	})
}

func Test_Service_IsRunning(t *testing.T) {
	t.Run("returns true when running", func(t *testing.T) {
		s := &service{
			semaphore: &semaphore{
				state: runningState,
			},
		}
		assert.True(t, s.GetStatus(context.Background()))
	})

	t.Run("returns false when stopped", func(t *testing.T) {
		s := &service{
			semaphore: &semaphore{
				state: stoppedState,
			},
		}
		assert.False(t, s.GetStatus(context.Background()))
	})
}
