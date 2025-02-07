package custom

import (
	"dillmann.com.br/nginx-ignition/core/certificate"
	"dillmann.com.br/nginx-ignition/core/common/dynamic_fields"
	"dillmann.com.br/nginx-ignition/core/common/validation"
)

type validationRules struct {
	dynamicFields []*dynamic_fields.DynamicField
}

func (r validationRules) DynamicFields() []*dynamic_fields.DynamicField {
	return r.dynamicFields
}

func (r validationRules) Validate(_ *certificate.IssueRequest) []validation.ConsistencyViolation {
	return make([]validation.ConsistencyViolation, 0)
}
