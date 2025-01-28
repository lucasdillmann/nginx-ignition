package scheduler

import (
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

func (s shutdown) Run() {
	log.Infof("Stopping scheduled tasks")
	s.scheduler.stop()
}
