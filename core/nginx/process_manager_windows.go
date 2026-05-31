//go:build windows

package nginx

import (
	"bytes"
	"errors"
	"os"
	"strings"
	"syscall"
	"time"

	"dillmann.com.br/nginx-ignition/core/common/log"
)

func (m *processManager) isPidAlive(pid int64) bool {
	_, err := os.FindProcess(int(pid))
	return err == nil
}

func (m *processManager) start() error {
	m.deleteTrafficStatsSocket()
	if err := m.runBackgroundCommand(time.Second * 2); err != nil {
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

	handle, err := syscall.OpenProcess(syscall.PROCESS_QUERY_INFORMATION, false, uint32(pid))
	if err != nil {
		return 0, err
	}

	defer syscall.CloseHandle(handle)
	var creationTime, exitTime, kernelTime, userTime syscall.Filetime

	if err := syscall.GetProcessTimes(handle, &creationTime, &exitTime, &kernelTime, &userTime); err != nil {
		return 0, err
	}

	startTime := time.Unix(0, creationTime.Nanoseconds())
	seconds := int64(time.Since(startTime).Seconds())
	if seconds < 0 {
		return 0, nil
	}

	return seconds, nil
}

func (m *processManager) runBackgroundCommand(waitDelay time.Duration, extraArgs ...string) error {
	cmd := m.prepareCommand(extraArgs...)

	var outputBuffer bytes.Buffer
	cmd.Stdout = &outputBuffer
	cmd.Stderr = &outputBuffer
	cmd.WaitDelay = waitDelay

	err := cmd.Start()
	if err != nil {
		return err
	}

	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case err := <-done:
		output := strings.TrimSpace(outputBuffer.String())
		if output != "" {
			return errors.New(output)
		}

		if err != nil {
			return err
		}

		return errors.New("nginx exited unexpectedly")
	case <-time.After(1 * time.Second):
		return cmd.Process.Release()
	}
}
