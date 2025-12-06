package healthcheck

import (
	"dillmann.com.br/nginx-ignition/core/common/healthcheck"
)

func toDto(status *healthcheck.Status) *statusDto {
	output := &statusDto{
		Healthy: status.Healthy,
		Details: make([]*detailDto, len(status.Details)),
	}

	for index, details := range status.Details {
		output.Details[index] = &detailDto{
			Component: details.ID,
			Healthy:   details.Error == nil,
		}
	}

	return output
}
