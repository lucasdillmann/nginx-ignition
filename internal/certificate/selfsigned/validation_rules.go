package selfsigned

import (
	"context"

	"github.com/lucasdillmann/nginx-ignition/internal/core/certificate"
	"github.com/lucasdillmann/nginx-ignition/internal/core/common/dynamicfields"
	"github.com/lucasdillmann/nginx-ignition/internal/core/common/validation"
)

type validationRules struct{}

func (r validationRules) DynamicFields() []dynamicfields.DynamicField {
	return make([]dynamicfields.DynamicField, 0)
}

func (r validationRules) Validate(
	_ context.Context,
	_ *certificate.IssueRequest,
) []validation.ConsistencyViolation {
	return make([]validation.ConsistencyViolation, 0)
}
