package scheduler

import (
	"context"

	"github.com/lucasdillmann/nginx-ignition/internal/core/common/lifecycle"
	"github.com/lucasdillmann/nginx-ignition/internal/core/common/log"
)

type shutdown struct {
	scheduler *Scheduler
}

func registerShutdown(lc *lifecycle.Lifecycle, sched *Scheduler) {
	lc.RegisterShutdown(shutdown{sched})
}

func (s shutdown) Priority() int {
	return shutdownPriority
}

func (s shutdown) Run(_ context.Context) {
	log.Infof("Stopping scheduled tasks")
	s.scheduler.stop()
}
