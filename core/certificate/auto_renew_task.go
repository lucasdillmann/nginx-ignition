package certificate

import (
	"context"
	"time"

	"dillmann.com.br/nginx-ignition/core/common/coreerror"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
	"dillmann.com.br/nginx-ignition/core/common/log"
	"dillmann.com.br/nginx-ignition/core/common/scheduler"
)

type autoRenewTask struct {
	service *service
}

func registerScheduledTask(
	ctx context.Context,
	service *service,
	sched *scheduler.Scheduler,
) error {
	task := autoRenewTask{service}
	return sched.Register(ctx, &task)
}

func (t autoRenewTask) Run(ctx context.Context) error {
	return t.service.renewAllDue(ctx)
}

func (t autoRenewTask) Schedule(ctx context.Context) (*scheduler.Schedule, error) {
	cfg, err := t.service.autoRenewSettings(ctx)
	if err != nil {
		return nil, err
	}

	var interval time.Duration

	switch cfg.IntervalUnit {
	case "MINUTES":
		interval = time.Minute * time.Duration(cfg.IntervalUnitCount)
	case "HOURS":
		interval = time.Hour * time.Duration(cfg.IntervalUnitCount)
	case "DAYS":
		interval = time.Hour * 24 * time.Duration(cfg.IntervalUnitCount)
	default:
		return nil, coreerror.New(i18n.M(ctx, i18n.K.CoreCertificateInvalidIntervalUnit), false)
	}

	return &scheduler.Schedule{
		Enabled:  cfg.Enabled,
		Interval: interval,
	}, nil
}

func (t autoRenewTask) OnScheduleStarted(ctx context.Context) {
	schedule, err := t.Schedule(ctx)
	if err != nil {
		return
	}

	log.Infof(
		"Certificate auto-renew task scheduled to run every %v minutes",
		schedule.Interval.Minutes(),
	)
}
