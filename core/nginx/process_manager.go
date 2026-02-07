package nginx

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"dillmann.com.br/nginx-ignition/core/common/configuration"
	"dillmann.com.br/nginx-ignition/core/common/log"
)

type processManager struct {
	binaryPath string
	configPath string
}

func newProcessManager(cfg *configuration.Configuration) (*processManager, error) {
	prefixedConfiguration := cfg.WithPrefix("nginx-ignition.nginx")
	binaryPath, err := prefixedConfiguration.Get("binary-path")
	if err != nil {
		return nil, err
	}

	configPath, err := prefixedConfiguration.Get("config-path")
	if err != nil {
		return nil, err
	}

	return &processManager{
		binaryPath: binaryPath,
		configPath: configPath,
	}, nil
}

func (m *processManager) sendReloadSignal() error {
	if err := m.runCommand("-s", "reload"); err != nil {
		return err
	}

	log.Infof("nginx reloaded")
	return nil
}

func (m *processManager) sendReopenSignal() error {
	log.Infof("Signaling nginx for log file reopen")
	return m.runCommand("-s", "reopen")
}

func (m *processManager) sendStopSignal() error {
	if err := m.runCommand("-s", "stop"); err != nil {
		return err
	}

	m.deleteTrafficStatsSocket()
	log.Infof("nginx stopped")
	return nil
}

func (m *processManager) deleteTrafficStatsSocket() {
	socketFile := filepath.Join(m.configPath, "traffic-stats.socket")
	if err := os.Remove(socketFile); err != nil && !os.IsNotExist(err) {
		log.Warnf("Failed to delete traffic-stats.socket: %v", err)
	}
}

func (m *processManager) currentPid() (int64, error) {
	pidFile := filepath.Join(m.configPath, "nginx.pid")
	data, err := os.ReadFile(pidFile)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, nil
		}
		return 0, err
	}

	pid, err := strconv.ParseInt(strings.TrimSpace(string(data)), 10, 64)
	if err != nil {
		return 0, err
	}

	if !m.isPidAlive(pid) {
		return 0, nil
	}

	return pid, nil
}

func (m *processManager) runCommand(extraArgs ...string) error {
	cmd := m.prepareCommand(extraArgs...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return errors.New(string(output))
	}

	return nil
}

func (m *processManager) prepareCommand(extraArgs ...string) *exec.Cmd {
	args := append(
		[]string{
			"-e", filepath.Join(m.configPath, "logs", "main.log"),
			"-c", filepath.Join(m.configPath, "config", "nginx.conf"),
		},
		extraArgs...,
	)

	return exec.Command(m.binaryPath, args...)
}
