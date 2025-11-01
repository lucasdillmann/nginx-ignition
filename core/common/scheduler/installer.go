package scheduler

import (
	"time"

	"dillmann.com.br/nginx-ignition/core/common/container"
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
