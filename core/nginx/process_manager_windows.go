//go:build windows

package nginx

import (
	"os"
)

func (m *processManager) isPidAlive(pid int64) bool {
	_, err := os.FindProcess(int(pid))
	return err == nil
}
