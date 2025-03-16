package certificate

import (
	"context"
	"dillmann.com.br/nginx-ignition/core/common/core_error"
	"dillmann.com.br/nginx-ignition/core/common/log"
	"dillmann.com.br/nginx-ignition/core/common/scheduler"
	"dillmann.com.br/nginx-ignition/core/settings"
	"time"
)

type autoRenewTask struct {
	service            *service
	settingsRepository *settings.Repository
}

func registerScheduledTask(
	ctx context.Context,
	service *service,
	settingsRepository settings.Repository,
	scheduler *scheduler.Scheduler,
) error {
	task := autoRenewTask{service, &settingsRepository}
	return scheduler.Register(ctx, &task)
}

func (t autoRenewTask) Run(ctx context.Context) error {
	return t.service.renewAllDue(ctx)
}

func (t autoRenewTask) Schedule(ctx context.Context) (*scheduler.Schedule, error) {
	cfg, err := (*t.settingsRepository).Get(ctx)
	if err != nil {
		return nil, err
	}

	var interval time.Duration

	certCfg := cfg.CertificateAutoRenew
	switch certCfg.IntervalUnit {
	case settings.MinutesTimeUnit:
		interval = time.Minute * time.Duration(certCfg.IntervalUnitCount)
	case settings.HoursTimeUnit:
		interval = time.Hour * time.Duration(certCfg.IntervalUnitCount)
	case settings.DaysTimeUnit:
		interval = time.Hour * 24 * time.Duration(certCfg.IntervalUnitCount)
	default:
		return nil, core_error.New("invalid interval unit", false)
	}

	return &scheduler.Schedule{
		Enabled:  cfg.CertificateAutoRenew.Enabled,
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
