package scheduler

import "time"

type Schedule struct {
	Enabled  bool
	Interval time.Duration
}

type Task interface {
	Run() error
	Schedule() (*Schedule, error)
	OnScheduleStarted()
}
