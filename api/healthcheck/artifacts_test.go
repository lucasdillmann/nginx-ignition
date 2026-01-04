package healthcheck

import (
	"dillmann.com.br/nginx-ignition/core/common/healthcheck"
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
