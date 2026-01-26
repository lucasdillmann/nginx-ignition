package validation

import (
	"dillmann.com.br/nginx-ignition/core/common/i18n"
)

type ConsistencyError struct {
	Violations []ConsistencyViolation
}

type ConsistencyViolation struct {
	Message *i18n.Message
	Path    string
}

func (err ConsistencyError) Error() string {
	return "One or more problems where found in the provided value"
}

func NewError(violations []ConsistencyViolation) *ConsistencyError {
	return &ConsistencyError{Violations: violations}
}

func SingleFieldError(path string, message *i18n.Message) *ConsistencyError {
	violation := ConsistencyViolation{
		Path:    path,
		Message: message,
	}

	return NewError([]ConsistencyViolation{violation})
}
