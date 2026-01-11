//go:build !windows

package nginx

import (
	"os"
	"syscall"

	"dillmann.com.br/nginx-ignition/core/common/log"
)

func (m *processManager) isPidAlive(pid int64) bool {
	process, err := os.FindProcess(int(pid))
	if err != nil {
		return false
	}

	return process.Signal(syscall.Signal(0)) == nil
}

func (m *processManager) start() error {
	if err := m.runCommand(); err != nil {
		return err
	}

	log.Infof("nginx started")
	return nil
}
