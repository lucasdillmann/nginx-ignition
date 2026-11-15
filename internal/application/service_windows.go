//go:build windows

package application

import (
	"context"
	"time"

	"github.com/lucasdillmann/nginx-ignition/internal/business/core/log"
	"golang.org/x/sys/windows/svc"
)

const (
	waitHint              = 30000
	statusUpdateInterval  = 10 * time.Second
	serviceExitCodeFailed = 1
)

type serviceHandler struct{}

func (h *serviceHandler) Execute(
	_ []string,
	requests <-chan svc.ChangeRequest,
	status chan<- svc.Status,
) (bool, uint32) {
	status <- svc.Status{
		State:      svc.StartPending,
		WaitHint:   waitHint,
		CheckPoint: 1,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ready := make(chan struct{})
	done := make(chan error, 1)

	go func() {
		done <- start(ctx, func() { close(ready) })
	}()

	ticker := time.NewTicker(statusUpdateInterval)
	defer ticker.Stop()

	checkpoint := uint32(1)

	for {
		select {
		case <-ready:
			status <- svc.Status{
				State:   svc.Running,
				Accepts: svc.AcceptStop | svc.AcceptShutdown,
			}

			return h.run(cancel, requests, status, done)

		case request := <-requests:
			switch request.Cmd {
			case svc.Interrogate:
				status <- request.CurrentStatus
			case svc.Stop, svc.Shutdown:
				return h.stop(cancel, requests, status, done)
			}

		case <-ticker.C:
			checkpoint++
			status <- svc.Status{
				State:      svc.StartPending,
				WaitHint:   waitHint,
				CheckPoint: checkpoint,
			}

		case err := <-done:
			if err != nil {
				log.Errorf("Service failed to start: %s", err)
				return true, serviceExitCodeFailed
			}

			return false, 0
		}
	}
}

func (h *serviceHandler) run(
	cancel context.CancelFunc,
	requests <-chan svc.ChangeRequest,
	status chan<- svc.Status,
	done <-chan error,
) (bool, uint32) {
	for {
		select {
		case request := <-requests:
			switch request.Cmd {
			case svc.Interrogate:
				status <- request.CurrentStatus
			case svc.Stop, svc.Shutdown:
				return h.stop(cancel, requests, status, done)
			}

		case err := <-done:
			if err != nil {
				log.Errorf("Service stopped unexpectedly: %s", err)
				return true, serviceExitCodeFailed
			}

			return false, 0
		}
	}
}

func (h *serviceHandler) stop(
	cancel context.CancelFunc,
	requests <-chan svc.ChangeRequest,
	status chan<- svc.Status,
	done <-chan error,
) (bool, uint32) {
	stopPending := svc.Status{
		State:      svc.StopPending,
		WaitHint:   waitHint,
		CheckPoint: 1,
	}
	status <- stopPending
	cancel()

	ticker := time.NewTicker(statusUpdateInterval)
	defer ticker.Stop()

	checkpoint := uint32(1)

	for {
		select {
		case request := <-requests:
			if request.Cmd == svc.Interrogate {
				status <- request.CurrentStatus
			}

		case <-ticker.C:
			checkpoint++
			stopPending.CheckPoint = checkpoint
			status <- stopPending

		case err := <-done:
			if err != nil {
				log.Errorf("Service stopped with error: %s", err)
				return true, serviceExitCodeFailed
			}

			return false, 0
		}
	}
}
