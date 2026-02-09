package nginx

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_processManager(t *testing.T) {
	tmpDir, _ := os.MkdirTemp("", "nginx-test")
	defer os.RemoveAll(tmpDir)

	manager := &processManager{
		configPath: tmpDir,
	}

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
}
