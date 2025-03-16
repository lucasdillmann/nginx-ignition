package nginx

import (
	"context"
	"dillmann.com.br/nginx-ignition/core/common/lifecycle"
	"dillmann.com.br/nginx-ignition/core/common/log"
)

type shutdown struct {
	command StopCommand
}

func registerShutdown(lifecycle *lifecycle.Lifecycle, command StopCommand) {
	lifecycle.RegisterShutdown(shutdown{command})
}

func (s shutdown) Priority() int {
	return shutdownPriority
}

func (s shutdown) Run(ctx context.Context) {
	if err := s.command(ctx, nil); err != nil {
		log.Warnf("Error stopping nginx: %s", err)
	}
}
