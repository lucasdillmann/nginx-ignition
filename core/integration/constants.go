package integration

import (
	"dillmann.com.br/nginx-ignition/core/common/core_error"
)

var (
	ErrIntegrationNotFound = core_error.New("Integration not found", true)
	ErrIntegrationDisabled = core_error.New("Integration is disabled", true)
)
