package scheduler

import (
	"time"

	"go.uber.org/dig"
)

func Install(container *dig.Container) error {
	if err := container.Provide(buildScheduler); err != nil {
		return err
	}

	if err := container.Invoke(registerStartup); err != nil {
		return err
	}

	return container.Invoke(registerShutdown)
}

func buildScheduler() *Scheduler {
	return &Scheduler{
		tickers: make(map[Task]*time.Ticker),
		stopped: false,
	}
}
