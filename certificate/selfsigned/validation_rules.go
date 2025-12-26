package selfsigned

import (
	"dillmann.com.br/nginx-ignition/core/certificate"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/validation"
)

type validationRules struct{}

func (r validationRules) DynamicFields() []dynamicfields.DynamicField {
	return make([]dynamicfields.DynamicField, 0)
}

func (r validationRules) Validate(_ *certificate.IssueRequest) []validation.ConsistencyViolation {
	return make([]validation.ConsistencyViolation, 0)
}
