package nginx

import (
	"context"
	"time"

	"dillmann.com.br/nginx-ignition/core/common/coreerror"
	"dillmann.com.br/nginx-ignition/core/common/log"
	"dillmann.com.br/nginx-ignition/core/common/scheduler"
	"dillmann.com.br/nginx-ignition/core/settings"
)

type logRotationTask struct {
	service          *service
	settingsCommands *settings.Commands
}

func registerScheduledTask(
	ctx context.Context,
	service *service,
	settingsCommands *settings.Commands,
	sched *scheduler.Scheduler,
) error {
	task := logRotationTask{service, settingsCommands}
	return sched.Register(ctx, &task)
}

func (t logRotationTask) Run(ctx context.Context) error {
	return t.service.rotateLogs(ctx)
}

func (t logRotationTask) Schedule(ctx context.Context) (*scheduler.Schedule, error) {
	cfg, err := t.settingsCommands.Get(ctx)
	if err != nil {
		return nil, err
	}

	var interval time.Duration

	certCfg := cfg.LogRotation
	switch certCfg.IntervalUnit {
	case settings.MinutesTimeUnit:
		interval = time.Minute * time.Duration(certCfg.IntervalUnitCount)
	case settings.HoursTimeUnit:
		interval = time.Hour * time.Duration(certCfg.IntervalUnitCount)
	case settings.DaysTimeUnit:
		interval = time.Hour * 24 * time.Duration(certCfg.IntervalUnitCount)
	default:
		return nil, coreerror.New("invalid interval unit", false)
	}

	return &scheduler.Schedule{
		Enabled:  cfg.LogRotation.Enabled,
		Interval: interval,
	}, nil
}

func (t logRotationTask) OnScheduleStarted(ctx context.Context) {
	schedule, err := t.Schedule(ctx)
	if err != nil {
		return
	}

	log.Infof(
		"Log rotation task scheduled to run every %v minutes",
		schedule.Interval.Minutes(),
	)
}
