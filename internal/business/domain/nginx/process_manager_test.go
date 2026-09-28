package nginx

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_processManager(t *testing.T) {
	tmpDir, _ := os.MkdirTemp("", "nginx-test")
	defer os.RemoveAll(tmpDir)

	manager := &processManager{
		configPath: tmpDir,
	}

	t.Run("ValidateConfiguration", func(t *testing.T) {
		if runtime.GOOS == "windows" {
			t.Skip("shell-based fake nginx is not available on Windows")
		}

		t.Run("uses the candidate configuration root", func(t *testing.T) {
			capturePath := filepath.Join(tmpDir, "args")
			fakeNginx := filepath.Join(tmpDir, "nginx-args")
			script := "#!/bin/sh\nprintf '%s\\n' \"$@\" > \"$CAPTURE_PATH\"\n"
			require.NoError(t, os.WriteFile(fakeNginx, []byte(script), 0o755))
			t.Setenv("CAPTURE_PATH", capturePath)

			candidateManager := &processManager{binaryPath: fakeNginx}
			candidateRoot := filepath.Join(tmpDir, "candidate")

			require.NoError(t, candidateManager.ValidateConfiguration(t.Context(), candidateRoot))

			contents, err := os.ReadFile(capturePath)
			require.NoError(t, err)
			assert.Equal(t, []string{
				"-e",
				filepath.Join(candidateRoot, "logs", "main.log"),
				"-c",
				filepath.Join(candidateRoot, "config", "nginx.conf"),
				"-t",
			}, strings.Split(strings.TrimSpace(string(contents)), "\n"))
		})

		t.Run("returns nginx test output", func(t *testing.T) {
			fakeNginx := filepath.Join(tmpDir, "nginx-failure")
			script := "#!/bin/sh\necho invalid configuration\nexit 1\n"
			require.NoError(t, os.WriteFile(fakeNginx, []byte(script), 0o755))

			failureManager := &processManager{binaryPath: fakeNginx}

			err := failureManager.ValidateConfiguration(t.Context(), tmpDir)

			assert.EqualError(t, err, "invalid configuration\n")
		})

		t.Run("returns context cancellation", func(t *testing.T) {
			fakeNginx := filepath.Join(tmpDir, "nginx-timeout")
			require.NoError(t, os.WriteFile(fakeNginx, []byte("#!/bin/sh\nsleep 5\n"), 0o755))

			timeoutManager := &processManager{binaryPath: fakeNginx}
			ctx, cancel := context.WithTimeout(t.Context(), 10*time.Millisecond)
			defer cancel()

			err := timeoutManager.ValidateConfiguration(ctx, tmpDir)

			assert.ErrorIs(t, err, context.DeadlineExceeded)
		})
	})

	t.Run("currentPid", func(t *testing.T) {
		t.Run("returns 0 when pid file does not exist", func(t *testing.T) {
			pid, err := manager.currentPid()
			assert.NoError(t, err)
			assert.Equal(t, int64(0), pid)
		})

		t.Run("returns 0 when pid is not alive", func(t *testing.T) {
			pidFile := filepath.Join(tmpDir, "nginx.pid")
			_ = os.WriteFile(pidFile, []byte("999999"), 0o644)

			pid, err := manager.currentPid()
			assert.NoError(t, err)
			assert.Equal(t, int64(0), pid)
		})
	})

	t.Run("deleteTrafficStatsSocket", func(t *testing.T) {
		t.Run("deletes file if it exists", func(t *testing.T) {
			socketFile := filepath.Join(tmpDir, "traffic-stats.socket")
			_ = os.WriteFile(socketFile, []byte("test"), 0o644)

			manager.deleteTrafficStatsSocket()

			assert.NoFileExists(t, socketFile)
		})

		t.Run("does nothing if file does not exist", func(t *testing.T) {
			socketFile := filepath.Join(tmpDir, "traffic-stats.socket")
			assert.NoFileExists(t, socketFile)

			manager.deleteTrafficStatsSocket()

			assert.NoFileExists(t, socketFile)
		})
	})

	t.Run("uptimeSeconds", func(t *testing.T) {
		pidFile := filepath.Join(tmpDir, "nginx.pid")
		_ = os.WriteFile(pidFile, []byte(strconv.Itoa(os.Getpid())), 0o644)

		uptime, err := manager.uptimeSeconds()
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, uptime, int64(0))
	})
}
