//go:build windows

package nginx

import (
	"bytes"
	"errors"
	"os"
	"strings"
	"time"

	"dillmann.com.br/nginx-ignition/core/common/log"
)

func (m *processManager) isPidAlive(pid int64) bool {
	_, err := os.FindProcess(int(pid))
	return err == nil
}

func (m *processManager) start() error {
	if err := m.runBackgroundCommand(time.Second * 2); err != nil {
		return err
	}

	log.Infof("nginx started")
	return nil
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
