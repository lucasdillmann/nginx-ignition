package commons

import (
	"dillmann.com.br/nginx-ignition/core/certificate/server"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/validation"
)

type DomainRules interface {
	DynamicFields() []*dynamicfields.DynamicField
	Validate(request *server.IssueRequest) []validation.ConsistencyViolation
}
