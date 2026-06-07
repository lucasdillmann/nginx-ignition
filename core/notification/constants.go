package notification

import "errors"

const (
	maxDeliveryAttempts = 5

	submissionDeliveryFailureLogFormat = "notification delivery failed for submission %s: %s"
)

var ErrConfigurationNotFound = errors.New("notification configuration not found")
