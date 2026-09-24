package nginx

import (
	"context"

	"github.com/lucasdillmann/nginx-ignition/internal/business/core/lifecycle"
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/log"
)

type shutdown struct {
	commands Commands
}

func registerShutdown(lc *lifecycle.Lifecycle, commands Commands) {
	lc.RegisterShutdown(shutdown{commands})
}

func (s shutdown) Priority() int {
	return shutdownPriority
}

func (s shutdown) Run(ctx context.Context) {
	if err := s.commands.Stop(ctx); err != nil {
		log.Warnf("Error stopping nginx: %s", err)
	}
}
