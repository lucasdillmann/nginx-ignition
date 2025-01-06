package access_list

import (
	"dillmann.com.br/nginx-ignition/core/common/validation"
)

type accessListValidator struct {
	errors []validation.ConsistencyViolation
}

func (validator *accessListValidator) validate(_ *AccessList) error {
	// TODO: Implement this

	if len(validator.errors) != 0 {
		return validation.NewError(validator.errors)
	} else {
		return nil
	}
}

func newValidator() *accessListValidator {
	return &accessListValidator{
		errors: make([]validation.ConsistencyViolation, 0),
	}
}
