package nginx

import (
	"dillmann.com.br/nginx-ignition/core/common/core_error"
	"dillmann.com.br/nginx-ignition/core/common/log"
	"dillmann.com.br/nginx-ignition/core/common/scheduler"
	"dillmann.com.br/nginx-ignition/core/settings"
	"time"
)

type logRotationTask struct {
	service            *service
	settingsRepository *settings.Repository
}

func registerScheduledTask(
	service *service,
	settingsRepository settings.Repository,
	scheduler *scheduler.Scheduler,
) error {
	task := logRotationTask{service, &settingsRepository}
	return scheduler.Register(&task)
}

func (t logRotationTask) Run() error {
	return t.service.rotateLogs()
}

func (t logRotationTask) Schedule() (*scheduler.Schedule, error) {
	cfg, err := (*t.settingsRepository).Get()
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
		return nil, core_error.New("invalid interval unit", false)
	}

	return &scheduler.Schedule{
		Enabled:  cfg.LogRotation.Enabled,
		Interval: interval,
	}, nil
}

func (t logRotationTask) OnScheduleStarted() {
	schedule, err := t.Schedule()
	if err != nil {
		return
	}

	log.Infof(
		"Log rotation task scheduled to run every %v minutes",
		schedule.Interval.Minutes(),
	)
}
