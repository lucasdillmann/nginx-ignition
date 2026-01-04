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
	"dillmann.com.br/nginx-ignition/core/host"
	"dillmann.com.br/nginx-ignition/core/settings"
)

func Test_logRotator(t *testing.T) {
	t.Run("readTail", func(t *testing.T) {
		rotator := &logRotator{}

		t.Run("returns all lines when count is less than max", func(t *testing.T) {
			tmpFile, _ := os.CreateTemp("", "test")
			defer os.Remove(tmpFile.Name())
			_, _ = tmpFile.WriteString("line1\nline2\n")
			_, _ = tmpFile.Seek(0, 0)

			lines, err := rotator.readTail(tmpFile, 5)
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

			lines, err := rotator.readTail(tmpFile, 2)
			assert.NoError(t, err)
			assert.Equal(t, []string{
				"line2",
				"line3",
			}, lines)
		})
	})

	t.Run("getLogFiles", func(t *testing.T) {
		ctx := context.Background()

		t.Run("returns main log and host logs", func(t *testing.T) {
			id1 := uuid.New()
			id2 := uuid.New()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := host.NewMockedCommands(ctrl)
			repo.EXPECT().GetAllEnabled(ctx).Return([]host.Host{
				{
					ID: id1,
				},
				{
					ID: id2,
				},
			}, nil)

			rotator := &logRotator{
				hostCommands: repo,
			}
			files, err := rotator.getLogFiles(ctx)

			assert.NoError(t, err)
			assert.Contains(t, files, "main.log")
			assert.Contains(t, files, fmt.Sprintf("host-%s.access.log", id1))
			assert.Contains(t, files, fmt.Sprintf("host-%s.error.log", id1))
			assert.Contains(t, files, fmt.Sprintf("host-%s.access.log", id2))
			assert.Contains(t, files, fmt.Sprintf("host-%s.error.log", id2))
		})
	})

	t.Run("rotate", func(t *testing.T) {
		ctx := context.Background()

		t.Run("rotates logs based on settings", func(t *testing.T) {
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
		})
	})
}
