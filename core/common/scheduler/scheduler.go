package scheduler

import (
	"context"
	"time"

	"dillmann.com.br/nginx-ignition/core/common/coreerror"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
	"dillmann.com.br/nginx-ignition/core/common/log"
)

var placeholderDuration = time.Hour * 99999

type Scheduler struct {
	tickers map[Task]*time.Ticker
	stopped bool
	started bool
}

func (s *Scheduler) Register(ctx context.Context, task Task) error {
	if s.stopped {
		return schedulerStoppedError(ctx)
	}

	s.tickers[task] = time.NewTicker(placeholderDuration)

	if s.started {
		return s.startTask(ctx, task)
	}

	return nil
}

func (s *Scheduler) start(ctx context.Context) error {
	if s.started {
		return coreerror.New(i18n.M(ctx, i18n.K.SchedulerErrorAlreadyStarted), false)
	}

	if s.stopped {
		return coreerror.New(i18n.M(ctx, i18n.K.SchedulerErrorShuttingDown), false)
	}

	s.started = true

	for task := range s.tickers {
		if err := s.startTask(ctx, task); err != nil {
			return err
		}
	}

	return nil
}

func (s *Scheduler) startTask(ctx context.Context, task Task) error {
	schedule, err := task.Schedule(ctx)
	if err != nil {
		return err
	}

	ticker := s.tickers[task]

	go func() {
		for range ticker.C {
			if !schedule.Enabled {
				return
			}

			if err := task.Run(ctx); err != nil {
				log.Errorf("Scheduled task failed with an error: %v", err)
			}
		}
	}()

	ticker.Reset(schedule.Interval)

	if schedule.Enabled {
		task.OnScheduleStarted(ctx)
	}

	return nil
}

func (s *Scheduler) stop() {
	s.stopped = true

	for task, ticker := range s.tickers {
		ticker.Stop()
		delete(s.tickers, task)
	}
}

func (s *Scheduler) Reload(ctx context.Context) error {
	if s.stopped {
		return schedulerStoppedError(ctx)
	}

	for task, ticker := range s.tickers {
		newSchedule, err := task.Schedule(ctx)
		if err != nil {
			return err
		}

		ticker.Reset(newSchedule.Interval)

		if newSchedule.Enabled {
			task.OnScheduleStarted(ctx)
		}
	}

	return nil
}

func schedulerStoppedError(ctx context.Context) error {
	return coreerror.New(i18n.M(ctx, i18n.K.SchedulerErrorShuttingDown), false)
}
