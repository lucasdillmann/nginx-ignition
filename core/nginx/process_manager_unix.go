//go:build !windows

package nginx

import (
	"errors"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"

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
	m.deleteTrafficStatsSocket()
	if err := m.runCommand(); err != nil {
		return err
	}

	log.Infof("nginx started")
	return nil
}

func (m *processManager) uptimeSeconds() (int64, error) {
	pid, err := m.currentPid()
	if err != nil {
		return 0, err
	}
	if pid == 0 {
		return 0, errors.New("nginx is not running")
	}

	return processUptimeSeconds(pid)
}

func processUptimeSeconds(pid int64) (int64, error) {
	command := exec.Command("ps", "-o", "lstart=", "-p", strconv.FormatInt(pid, 10))
	output, err := command.Output()
	if err != nil {
		return 0, err
	}

	startText := strings.TrimSpace(string(output))
	if startText == "" {
		return 0, errors.New("process not found")
	}

	startTime, err := time.ParseInLocation("Mon Jan _2 15:04:05 2006", startText, time.Local)
	if err != nil {
		return 0, err
	}

	seconds := int64(time.Since(startTime).Seconds())
	if seconds < 0 {
		return 0, nil
	}
	return seconds, nil
}
