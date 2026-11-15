//go:build windows

package application

import (
	"context"
	"os"
	"os/signal"

	"github.com/lucasdillmann/nginx-ignition/internal/business/core/constants"
	"golang.org/x/sys/windows/svc"
)

func Run() error {
	isService, err := svc.IsWindowsService()
	if err != nil {
		return err
	}

	if isService {
		return svc.Run(constants.ServiceName, &serviceHandler{})
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	return start(ctx, nil)
}
