package application

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lucasdillmann/nginx-ignition/internal/core/common/container"
	"github.com/lucasdillmann/nginx-ignition/internal/core/common/lifecycle"
	"github.com/lucasdillmann/nginx-ignition/internal/core/common/log"
)

func Start() error {
	startTime := time.Now().UnixNano() / int64(time.Millisecond)
	ctx := context.Background()

	if err := startContainer(ctx); err != nil {
		return err
	}

	return container.Run(func(lifecycle *lifecycle.Lifecycle) error {
		return runLifecycle(ctx, lifecycle, startTime)
	})
}

func runLifecycle(ctx context.Context, lc *lifecycle.Lifecycle, startTime int64) error {
	if err := lc.FireStartup(ctx); err != nil {
		return err
	}

	endTime := time.Now().UnixNano() / int64(time.Millisecond)
	log.Infof("Application started in %d ms", endTime-startTime)

	receivedSignal := waitForShutdownSignal()

	log.Infof(
		"Application shutdown signal received (%s). Starting graceful shutdown.",
		receivedSignal,
	)
	lc.FireShutdown(ctx)

	log.Infof("Shutdown complete")
	return nil
}

func waitForShutdownSignal() os.Signal {
	channel := make(chan os.Signal, 1)
	signal.Notify(
		channel,
		syscall.SIGHUP,
		syscall.SIGQUIT,
		syscall.SIGTERM,
		syscall.SIGINT,
		syscall.SIGSEGV,
	)

	return <-channel
}
