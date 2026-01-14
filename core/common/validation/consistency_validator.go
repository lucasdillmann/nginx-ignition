package validation

import (
	"dillmann.com.br/nginx-ignition/i18n"
)

type ConsistencyValidator struct {
	violations []ConsistencyViolation
}

func NewValidator() *ConsistencyValidator {
	return &ConsistencyValidator{}
}

func (v *ConsistencyValidator) Add(path string, message *i18n.Message) {
	v.violations = append(v.violations, ConsistencyViolation{
		Path:    path,
		Message: message,
	})
}

func (v *ConsistencyValidator) Result() error {
	if len(v.violations) == 0 {
		return nil
	}

	return NewError(v.violations)
}
