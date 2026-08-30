package healthcheck

import (
	"github.com/lucasdillmann/nginx-ignition/internal/core/common/healthcheck"
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
