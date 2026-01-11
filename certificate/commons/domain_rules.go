package commons

import (
	"context"

	"dillmann.com.br/nginx-ignition/core/certificate"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/validation"
)

type DomainRules interface {
	DynamicFields() []dynamicfields.DynamicField
	Validate(
		ctx context.Context,
		request *certificate.IssueRequest,
	) []validation.ConsistencyViolation
}
