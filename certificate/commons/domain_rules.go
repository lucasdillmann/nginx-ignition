package commons

import (
	"dillmann.com.br/nginx-ignition/core/certificate"
	"dillmann.com.br/nginx-ignition/core/common/dynamic_fields"
	"dillmann.com.br/nginx-ignition/core/common/validation"
)

type DomainRules interface {
	DynamicFields() []*dynamic_fields.DynamicField
	Validate(request *certificate.IssueRequest) []validation.ConsistencyViolation
}
