package scheduler

import (
	"context"
	"time"
)

type Schedule struct {
	Enabled  bool
	Interval time.Duration
}

type Task interface {
	Run(ctx context.Context) error
	Schedule(ctx context.Context) (*Schedule, error)
	OnScheduleStarted(ctx context.Context)
}
