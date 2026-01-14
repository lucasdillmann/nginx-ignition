package coreerror

import (
	"dillmann.com.br/nginx-ignition/i18n"
)

type CoreError struct {
	Message     *i18n.Message
	UserRelated bool
}

func New(message *i18n.Message, blameUser bool) *CoreError {
	return &CoreError{
		Message:     message,
		UserRelated: blameUser,
	}
}

func (e CoreError) Error() string {
	return e.Message.String()
}
