package notification

import (
	"context"
	"time"

	"dillmann.com.br/nginx-ignition/core/common/log"
	"dillmann.com.br/nginx-ignition/core/common/scheduler"
)

const deliveryTaskInterval = 30 * time.Second

type deliveryTask struct {
	commands Commands
}

func registerScheduledTask(
	ctx context.Context,
	commands Commands,
	sched *scheduler.Scheduler,
) error {
	task := deliveryTask{commands: commands}
	return sched.Register(ctx, &task)
}

func (t *deliveryTask) Run(ctx context.Context) error {
	return t.commands.ProcessPendingDeliveries(ctx)
}

func (t *deliveryTask) Schedule(_ context.Context) (*scheduler.Schedule, error) {
	return &scheduler.Schedule{
		Enabled:  true,
		Interval: deliveryTaskInterval,
	}, nil
}

func (t *deliveryTask) OnScheduleStarted(_ context.Context) {
	log.Infof(
		"Notification delivery task scheduled to run every %v seconds",
		deliveryTaskInterval.Seconds(),
	)
}
