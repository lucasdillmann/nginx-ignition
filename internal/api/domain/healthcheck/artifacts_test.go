package healthcheck

import (
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/healthcheck"
)

func newHealthcheckStatus() *healthcheck.Status {
	return &healthcheck.Status{
		Healthy: true,
		Details: []healthcheck.Detail{
			{
				ID:    "db",
				Error: nil,
			},
			{
				ID:    "nginx",
				Error: nil,
			},
		},
	}
}
