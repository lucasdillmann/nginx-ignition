package healthcheck

import (
	"context"
)

type HealthCheck struct {
	providers []Provider
}

func (h *HealthCheck) Register(provider Provider) {
	h.providers = append(h.providers, provider)
}

func (h *HealthCheck) Status(ctx context.Context) *Status {
	status := Status{
		Healthy: true,
		Details: make([]*Detail, len(h.providers)),
	}

	for index, provider := range h.providers {
		err := provider.Check(ctx)
		if err != nil {
			status.Healthy = false
		}

		status.Details[index] = &Detail{
			ID:    provider.ID(),
			Error: err,
		}
	}

	return &status
}
