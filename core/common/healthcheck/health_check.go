package healthcheck

import (
	"context"
)

type HealthCheck struct {
	providers []Provider
}

func New() *HealthCheck {
	return &HealthCheck{
		providers: make([]Provider, 0),
	}
}

func (h *HealthCheck) Register(provider Provider) {
	h.providers = append(h.providers, provider)
}

func (h *HealthCheck) Status(ctx context.Context) *Status {
	status := Status{
		Healthy: true,
		Details: make([]Detail, len(h.providers)),
	}

	for index, provider := range h.providers {
		status.Details[index] = Detail{
			ID:    provider.ID(),
			Error: provider.Check(ctx),
		}

		if status.Details[index].Error != nil {
			status.Healthy = false
		}
	}

	return &status
}
