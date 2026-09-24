package domainrules

import (
	"context"

	"github.com/lucasdillmann/nginx-ignition/internal/business/core/dynamicfields"
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/validation"
	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/certificate"
)

type DomainRules interface {
	DynamicFields() []dynamicfields.DynamicField
	Validate(
		ctx context.Context,
		request *certificate.IssueRequest,
	) []validation.ConsistencyViolation
}
