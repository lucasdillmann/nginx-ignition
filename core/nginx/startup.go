package nginx

import (
	"context"
	"dillmann.com.br/nginx-ignition/core/common/lifecycle"
)

type startup struct {
	service *service
}

func registerStartup(lifecycle *lifecycle.Lifecycle, service *service) {
	lifecycle.RegisterStartup(startup{service})
}

func (s startup) Run(ctx context.Context) error {
	go func() {
		s.service.attachListeners()
	}()

	return s.service.start(ctx)
}

func (s startup) Priority() int {
	return startupPriority
}

func (s startup) Async() bool {
	return true
}
