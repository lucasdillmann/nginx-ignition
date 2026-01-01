package nginx

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"dillmann.com.br/nginx-ignition/core/common/configuration"
)

func TestLogReader_Read(t *testing.T) {
	ctx := context.Background()
	tmpDir, _ := os.MkdirTemp("", "logs")
	defer os.RemoveAll(tmpDir)

	_ = os.Mkdir(filepath.Join(tmpDir, "logs"), 0o755)
	logFile := filepath.Join(tmpDir, "logs", "test.log")
	_ = os.WriteFile(logFile, []byte("line1\nline2\nline3\n"), 0o644)

	t.Setenv("NGINX_IGNITION_NGINX_CONFIG_PATH", tmpDir)

	cfg := configuration.New()
	reader := newLogReader(cfg)

	t.Run("reads and reverses lines correctly", func(t *testing.T) {
		lines, err := reader.read(ctx, "test.log", 10)
		assert.NoError(t, err)
		assert.Equal(t, []string{
			"line3",
			"line2",
			"line1",
		}, lines)
	})

	t.Run("tails and reverses lines correctly", func(t *testing.T) {
		lines, err := reader.read(ctx, "test.log", 2)
		assert.NoError(t, err)
		assert.Equal(t, []string{
			"line3",
			"line2",
		}, lines)
	})
}
