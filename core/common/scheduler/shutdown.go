package scheduler

import (
	"context"

	"dillmann.com.br/nginx-ignition/core/common/lifecycle"
	"dillmann.com.br/nginx-ignition/core/common/log"
)

type shutdown struct {
	scheduler *Scheduler
}

func registerShutdown(lifecycle *lifecycle.Lifecycle, scheduler *Scheduler) {
	lifecycle.RegisterShutdown(shutdown{scheduler})
}

func (s shutdown) Priority() int {
	return shutdownPriority
}

func (s shutdown) Run(_ context.Context) {
	log.Infof("Stopping scheduled tasks")
	s.scheduler.stop()
}
