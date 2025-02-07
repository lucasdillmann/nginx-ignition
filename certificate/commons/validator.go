package commons

import (
	"dillmann.com.br/nginx-ignition/core/certificate"
	"dillmann.com.br/nginx-ignition/core/common/constants"
	"dillmann.com.br/nginx-ignition/core/common/dynamic_fields"
	"dillmann.com.br/nginx-ignition/core/common/validation"
	"fmt"
)

func Validate(request *certificate.IssueRequest, domainRules DomainRules) error {
	violations := append(validateBaseFields(request), domainRules.Validate(request)...)

	dynamicFieldsResult := dynamic_fields.Validate(domainRules.DynamicFields(), request.Parameters)
	if dynamicFieldsResult != nil {
		violations = append(violations, dynamicFieldsResult.Violations...)
	}

	if len(violations) > 0 {
		return validation.NewError(violations)
	}

	return nil
}

func validateBaseFields(request *certificate.IssueRequest) []validation.ConsistencyViolation {
	violations := make([]validation.ConsistencyViolation, 0)
	if len(request.DomainNames) == 0 {
		violations = append(violations, validation.ConsistencyViolation{
			Path:    "domainNames",
			Message: "At least one domain name must be informed",
		})
	}

	for index, domainName := range request.DomainNames {
		if !constants.TLDPattern.MatchString(domainName) {
			violations = append(violations, validation.ConsistencyViolation{
				Path:    fmt.Sprintf("domainNames[%d]", index),
				Message: "Value is not a valid domain name",
			})
		}
	}

	return violations
}
