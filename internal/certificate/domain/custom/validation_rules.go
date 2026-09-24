package custom

import (
	"context"

	"github.com/lucasdillmann/nginx-ignition/internal/business/core/dynamicfields"
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/validation"
	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/certificate"
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
