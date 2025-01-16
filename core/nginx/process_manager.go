package nginx

import (
	"dillmann.com.br/nginx-ignition/core/common/configuration"
	"dillmann.com.br/nginx-ignition/core/common/log"
	"errors"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type processManager struct {
	binaryPath string
	configPath string
}

func newProcessManager(configuration *configuration.Configuration) (*processManager, error) {
	prefixedConfiguration := configuration.WithPrefix("nginx-ignition.nginx")
	binaryPath, err := prefixedConfiguration.Get("binary-path")
	if err != nil {
		return nil, err
	}

	configPath, err := prefixedConfiguration.Get("config-directory")
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

	log.Infof("nginx stopped")
	return nil
}

func (m *processManager) start() error {
	if err := m.runCommand(); err != nil {
		return err
	}

	log.Infof("nginx started")
	return nil
}

func (m *processManager) currentPid() (int64, error) {
	pidFile := m.configPath + "/nginx.pid"
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

func (m *processManager) isPidAlive(pid int64) bool {
	cmd := exec.Command("kill", "-0", strconv.FormatInt(pid, 10))
	err := cmd.Run()
	return err == nil
}

func (m *processManager) runCommand(extraArgs ...string) error {
	args := append([]string{m.binaryPath, "-c", m.configPath + "/config/nginx.conf"}, extraArgs...)

	cmd := exec.Command(args[0], args[1:]...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return errors.New(string(output))
	}

	return nil
}
