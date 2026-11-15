//go:build !windows

package application

import (
	"context"
	"os/signal"
	"syscall"
)

func Run() error {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGHUP,
		syscall.SIGQUIT,
		syscall.SIGTERM,
		syscall.SIGINT,
	)
	defer stop()

	return start(ctx, nil)
}
