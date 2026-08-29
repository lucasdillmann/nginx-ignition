package custom

import (
	"context"

	"nginx-ignition/internal/core/certificate"
	"nginx-ignition/internal/core/common/dynamicfields"
	"nginx-ignition/internal/core/common/validation"
)

type validationRules struct {
	dynamicFields []dynamicfields.DynamicField
}

func (r validationRules) DynamicFields() []dynamicfields.DynamicField {
	return r.dynamicFields
}

func (r validationRules) Validate(
	_ context.Context,
	_ *certificate.IssueRequest,
) []validation.ConsistencyViolation {
	return make([]validation.ConsistencyViolation, 0)
}
