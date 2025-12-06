package healthcheck

import (
	"dillmann.com.br/nginx-ignition/core/common/container"
)

func Install() error {
	return container.Provide(newHealthCheck)
}

func newHealthCheck() *HealthCheck {
	return &HealthCheck{
		providers: make([]Provider, 0),
	}
}
