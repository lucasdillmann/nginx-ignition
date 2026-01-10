//go:build !windows

package nginx

import (
	"os/exec"
	"strconv"
)

func (m *processManager) isPidAlive(pid int64) bool {
	cmd := exec.Command("kill", "-0", strconv.FormatInt(pid, 10))
	return cmd.Run() == nil
}
