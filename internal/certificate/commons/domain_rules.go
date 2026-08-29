package commons

import (
	"context"

	"github.com/lucasdillmann/nginx-ignition/internal/core/certificate"
	"github.com/lucasdillmann/nginx-ignition/internal/core/common/dynamicfields"
	"github.com/lucasdillmann/nginx-ignition/internal/core/common/validation"
)

type DomainRules interface {
	DynamicFields() []dynamicfields.DynamicField
	Validate(
		ctx context.Context,
		request *certificate.IssueRequest,
	) []validation.ConsistencyViolation
}
