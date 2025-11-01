package boot

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"dillmann.com.br/nginx-ignition/core/common/container"
	"dillmann.com.br/nginx-ignition/core/common/lifecycle"
	"dillmann.com.br/nginx-ignition/core/common/log"
)

func StartApplication() error {
	startTime := time.Now().UnixNano() / int64(time.Millisecond)
	ctx := context.Background()

	if err := startContainer(ctx); err != nil {
		return err
	}

	return container.Run(func(lifecycle *lifecycle.Lifecycle) error {
		return runLifecycle(ctx, lifecycle, startTime)
	})
}

func runLifecycle(ctx context.Context, lifecycle *lifecycle.Lifecycle, startTime int64) error {
	if err := lifecycle.FireStartup(ctx); err != nil {
		return err
	}

	endTime := time.Now().UnixNano() / int64(time.Millisecond)
	log.Infof("Application started in %d ms", endTime-startTime)

	receivedSignal := waitForShutdownSignal()

	log.Infof("Application shutdown signal received (%s). Starting graceful shutdown.", receivedSignal)
	lifecycle.FireShutdown(ctx)

	log.Infof("Shutdown complete")
	return nil
}

func waitForShutdownSignal() os.Signal {
	channel := make(chan os.Signal, 1)
	signal.Notify(channel, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGSEGV)

	return <-channel
}
