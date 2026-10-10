package authorization

import (
	"context"

	"github.com/lucasdillmann/nginx-ignition/internal/business/core/lifecycle"
)

type startup struct {
	service *service
}

func registerStartup(lc *lifecycle.Lifecycle, service *service) {
	commandInstance := &startup{service}
	lc.RegisterStartup(commandInstance)
}

func (s startup) Run(ctx context.Context) error {
	return s.service.initialize(ctx)
}

func (s startup) Priority() int {
	return startupPriority
}

func (s startup) Async() bool {
	return false
}
