package nginx

import (
	"context"
	"time"

	"github.com/lucasdillmann/nginx-ignition/internal/business/core/coreerror"
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/i18n"
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/log"
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/scheduler"
	settings2 "github.com/lucasdillmann/nginx-ignition/internal/business/domain/settings"
)

type logRotationTask struct {
	service          *service
	settingsCommands settings2.Commands
}

func registerScheduledTask(
	ctx context.Context,
	service *service,
	settingsCommands settings2.Commands,
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
	case settings2.MinutesTimeUnit:
		interval = time.Minute * time.Duration(certCfg.IntervalUnitCount)
	case settings2.HoursTimeUnit:
		interval = time.Hour * time.Duration(certCfg.IntervalUnitCount)
	case settings2.DaysTimeUnit:
		interval = time.Hour * 24 * time.Duration(certCfg.IntervalUnitCount)
	default:
		return nil, coreerror.New(i18n.M(ctx, i18n.K.CommonInvalidIntervalUnit), false)
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
