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

	"dillmann.com.br/nginx-ignition/core/common/configuration"
	"dillmann.com.br/nginx-ignition/core/host"
	"dillmann.com.br/nginx-ignition/core/settings"
)

func TestLogRotator_ReadTail(t *testing.T) {
	r := &logRotator{}

	t.Run("returns all lines when count is less than max", func(t *testing.T) {
		tmpFile, _ := os.CreateTemp("", "test")
		defer os.Remove(tmpFile.Name())
		_, _ = tmpFile.WriteString("line1\nline2\n")
		_, _ = tmpFile.Seek(0, 0)

		lines, err := r.readTail(tmpFile, 5)
		assert.NoError(t, err)
		assert.Equal(t, []string{
			"line1",
			"line2",
		}, lines)
	})

	t.Run("trims lines when count exceeds max", func(t *testing.T) {
		tmpFile, _ := os.CreateTemp("", "test")
		defer os.Remove(tmpFile.Name())
		_, _ = tmpFile.WriteString("line1\nline2\nline3\n")
		_, _ = tmpFile.Seek(0, 0)

		lines, err := r.readTail(tmpFile, 2)
		assert.NoError(t, err)
		assert.Equal(t, []string{
			"line2",
			"line3",
		}, lines)
	})
}

func TestLogRotator_GetLogFiles(t *testing.T) {
	ctx := context.Background()

	t.Run("returns main log and host logs", func(t *testing.T) {
		id1 := uuid.New()
		id2 := uuid.New()

		repo := &host.Commands{
			GetAllEnabled: func(_ context.Context) ([]host.Host, error) {
				return []host.Host{
					{
						ID: id1,
					},
					{
						ID: id2,
					},
				}, nil
			},
		}

		r := &logRotator{
			hostCommands: repo,
		}
		files, err := r.getLogFiles(ctx)

		assert.NoError(t, err)
		assert.Contains(t, files, "main.log")
		assert.Contains(t, files, fmt.Sprintf("host-%s.access.log", id1))
		assert.Contains(t, files, fmt.Sprintf("host-%s.error.log", id1))
		assert.Contains(t, files, fmt.Sprintf("host-%s.access.log", id2))
		assert.Contains(t, files, fmt.Sprintf("host-%s.error.log", id2))
	})
}

func TestLogRotator_Rotate(t *testing.T) {
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

	settingsCmds := &settings.Commands{
		Get: func(_ context.Context) (*settings.Settings, error) {
			return &settings.Settings{
				LogRotation: &settings.LogRotationSettings{
					Enabled:      true,
					MaximumLines: 2,
				},
			}, nil
		},
	}

	hostCmds := &host.Commands{
		GetAllEnabled: func(_ context.Context) ([]host.Host, error) {
			return []host.Host{}, nil
		},
	}

	pm := &processManager{
		binaryPath: fakeNginx,
		configPath: tmpDir,
	}

	rotator := newLogRotator(cfg, settingsCmds, hostCmds, pm)

	err = rotator.rotate(ctx)
	assert.NoError(t, err)

	content, err := os.ReadFile(mainLogPath)
	require.NoError(t, err)
	assert.Equal(t, "line2\nline3\n", string(content))
}
