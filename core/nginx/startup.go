package nginx

import "dillmann.com.br/nginx-ignition/core/common/lifecycle"

type startup struct {
	service *service
}

func registerStartup(lifecycle *lifecycle.Lifecycle, service *service) {
	lifecycle.RegisterStartup(startup{service})
}

func (s startup) Run() error {
	go func() {
		s.service.attachListeners()
	}()

	return s.service.start()
}

func (s startup) Priority() int {
	return startupPriority
}

func (s startup) Async() bool {
	return true
}
