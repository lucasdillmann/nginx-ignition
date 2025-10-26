package scheduler

import (
	"context"
	"time"

	"dillmann.com.br/nginx-ignition/core/common/core_error"
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
		return schedulerStoppedError()
	}

	s.tickers[task] = time.NewTicker(placeholderDuration)

	if s.started {
		return s.startTask(ctx, task)
	}

	return nil
}

func (s *Scheduler) start(ctx context.Context) error {
	if s.started {
		return core_error.New("Scheduler already started", false)
	}

	if s.stopped {
		return core_error.New("Scheduler is shutting-down or was already stopped", false)
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
		return schedulerStoppedError()
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

func schedulerStoppedError() error {
	return core_error.New("Scheduler is shutting-down or was already stopped", false)
}
