package integration

import (
	"dillmann.com.br/nginx-ignition/core/common/coreerror"
)

var (
	ErrIntegrationNotFound = coreerror.New("Integration not found", true)
	ErrIntegrationDisabled = coreerror.New("Integration is disabled", true)
)
