package healthcheck

import (
	"dillmann.com.br/nginx-ignition/core/common/healthcheck"
	"dillmann.com.br/nginx-ignition/core/common/ptr"
)

func toDto(status *healthcheck.Status) *statusDto {
	output := &statusDto{
		Healthy: status.Healthy,
		Details: make([]*detailDto, len(status.Details)),
	}

	for index, details := range status.Details {
		var reason *string

		if details.Error != nil {
			reason = ptr.Of(details.Error.Error())
		}

		output.Details[index] = &detailDto{
			Component: details.ID,
			Healthy:   details.Healthy,
			Reason:    reason,
		}
	}

	return output
}
