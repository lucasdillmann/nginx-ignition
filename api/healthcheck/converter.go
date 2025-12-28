package healthcheck

import (
	"dillmann.com.br/nginx-ignition/core/common/healthcheck"
)

func toDTO(status *healthcheck.Status) *statusDTO {
	output := &statusDTO{
		Healthy: status.Healthy,
		Details: make([]detailDTO, len(status.Details)),
	}

	for index, details := range status.Details {
		output.Details[index] = detailDTO{
			Component: details.ID,
			Healthy:   details.Error == nil,
		}
	}

	return output
}
