package user

import (
	"context"
	"time"

	"github.com/lucasdillmann/nginx-ignition/internal/business/core/scheduler"
)

const sessionCleanupInterval = time.Minute * 5

type sessionCleanupTask struct {
	service *service
}

func registerSessionCleanupTask(
	ctx context.Context,
	service *service,
	sched *scheduler.Scheduler,
) error {
	task := sessionCleanupTask{service}
	return sched.Register(ctx, &task)
}

func (t sessionCleanupTask) Run(ctx context.Context) error {
	return t.service.deleteExpiredSessions(ctx)
}

func (t sessionCleanupTask) Schedule(_ context.Context) (*scheduler.Schedule, error) {
	return &scheduler.Schedule{
		Enabled:  true,
		Interval: sessionCleanupInterval,
	}, nil
}

func (t sessionCleanupTask) OnScheduleStarted(_ context.Context) {
	// NO-OP
}
