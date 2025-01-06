package user

import (
	"dillmann.com.br/nginx-ignition/core/common/core_errors"
	"dillmann.com.br/nginx-ignition/core/common/validation"
)

type validator struct {
	errors []validation.ConsistencyViolation
}

func (v *validator) validate(_ *User) error {
	return core_errors.NotImplemented()
}

func newValidator() *validator {
	return &validator{
		errors: []validation.ConsistencyViolation{},
	}
}
