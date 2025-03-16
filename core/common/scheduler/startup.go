package scheduler

import (
	"context"
	"dillmann.com.br/nginx-ignition/core/common/lifecycle"
	"dillmann.com.br/nginx-ignition/core/common/log"
)

type startup struct {
	scheduler *Scheduler
}

func registerStartup(lifecycle *lifecycle.Lifecycle, scheduler *Scheduler) {
	lifecycle.RegisterStartup(startup{scheduler})
}

func (s startup) Run(ctx context.Context) error {
	log.Infof("Starting scheduled tasks")
	return s.scheduler.start(ctx)
}

func (s startup) Priority() int {
	return startupPriority
}

func (s startup) Async() bool {
	return false
}
