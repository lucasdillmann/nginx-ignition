package notification

import "errors"

const maxDeliveryAttempts = 5

var ErrConfigurationNotFound = errors.New("notification configuration not found")
