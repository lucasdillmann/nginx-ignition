package nginx

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ProcessManager(t *testing.T) {
	tmpDir, _ := os.MkdirTemp("", "nginx-test")
	defer os.RemoveAll(tmpDir)

	m := &processManager{
		configPath: tmpDir,
	}

	t.Run("currentPid", func(t *testing.T) {
		t.Run("returns 0 when pid file does not exist", func(t *testing.T) {
			pid, err := m.currentPid()
			assert.NoError(t, err)
			assert.Equal(t, int64(0), pid)
		})

		t.Run("returns 0 when pid is not alive", func(t *testing.T) {
			pidFile := filepath.Join(tmpDir, "nginx.pid")
			_ = os.WriteFile(pidFile, []byte("999999"), 0o644)

			pid, err := m.currentPid()
			assert.NoError(t, err)
			assert.Equal(t, int64(0), pid)
		})
	})
}
