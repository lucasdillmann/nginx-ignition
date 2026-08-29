package commons

import (
	"context"

	"nginx-ignition/internal/core/certificate"
	"nginx-ignition/internal/core/common/dynamicfields"
	"nginx-ignition/internal/core/common/validation"
)

type DomainRules interface {
	DynamicFields() []dynamicfields.DynamicField
	Validate(
		ctx context.Context,
		request *certificate.IssueRequest,
	) []validation.ConsistencyViolation
}
