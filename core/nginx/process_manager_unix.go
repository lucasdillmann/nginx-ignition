//go:build !windows

package nginx

import (
	"os/exec"
	"strconv"

	"dillmann.com.br/nginx-ignition/core/common/log"
)

func (m *processManager) isPidAlive(pid int64) bool {
	cmd := exec.Command("kill", "-0", strconv.FormatInt(pid, 10))
	return cmd.Run() == nil
}

func (m *processManager) start() error {
	if err := m.runCommand(); err != nil {
		return err
	}

	log.Infof("nginx started")
	return nil
}
