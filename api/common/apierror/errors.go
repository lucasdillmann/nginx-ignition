package apierror

import (
	"dillmann.com.br/nginx-ignition/core/common/i18n"
)

type APIError struct {
	Message    *i18n.Message
	StatusCode int
}

func New(statusCode int, message *i18n.Message) *APIError {
	return &APIError{
		Message:    message,
		StatusCode: statusCode,
	}
}

func (e APIError) Error() string {
	return e.Message.String()
}
