package application

import (
	"context"
	"time"

	"github.com/lucasdillmann/nginx-ignition/internal/business/core/container"
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/lifecycle"
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/log"
)

func start(ctx context.Context, onReady func()) error {
	startTime := time.Now().UnixNano() / int64(time.Millisecond)

	if err := startContainer(ctx); err != nil {
		return err
	}

	return container.Run(func(lc *lifecycle.Lifecycle) error {
		return runLifecycle(ctx, lc, startTime, onReady)
	})
}

func runLifecycle(
	ctx context.Context,
	lc *lifecycle.Lifecycle,
	startTime int64,
	onReady func(),
) error {
	if err := lc.FireStartup(ctx); err != nil {
		return err
	}

	endTime := time.Now().UnixNano() / int64(time.Millisecond)
	log.Infof("Application started in %d ms", endTime-startTime)

	if onReady != nil {
		onReady()
	}

	<-ctx.Done()

	log.Info("Application shutdown requested. Starting graceful shutdown.")
	lc.FireShutdown(context.WithoutCancel(ctx))
	log.Info("Shutdown complete")

	return nil
}
