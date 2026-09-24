package scheduler

import (
	"time"

	"github.com/lucasdillmann/nginx-ignition/internal/business/core/container"
)

func Install() error {
	if err := container.Provide(buildScheduler); err != nil {
		return err
	}

	return container.Run(registerStartup, registerShutdown)
}

func buildScheduler() *Scheduler {
	return &Scheduler{
		tickers: make(map[Task]*time.Ticker),
		stopped: false,
	}
}
