package commons

import (
	"context"
	"fmt"

	"dillmann.com.br/nginx-ignition/core/certificate"
	"dillmann.com.br/nginx-ignition/core/common/constants"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/validation"
	"dillmann.com.br/nginx-ignition/i18n"
)

func Validate(
	ctx context.Context,
	request *certificate.IssueRequest,
	domainRules DomainRules,
) error {
	violations := append(validateBaseFields(ctx, request), domainRules.Validate(ctx, request)...)

	dynamicFieldsResult := dynamicfields.Validate(
		ctx,
		domainRules.DynamicFields(),
		request.Parameters,
	)
	if dynamicFieldsResult != nil {
		violations = append(violations, dynamicFieldsResult.Violations...)
	}

	if len(violations) > 0 {
		return validation.NewError(violations)
	}

	return nil
}

func validateBaseFields(
	ctx context.Context,
	request *certificate.IssueRequest,
) []validation.ConsistencyViolation {
	violations := make([]validation.ConsistencyViolation, 0)
	if len(request.DomainNames) == 0 {
		violations = append(violations, validation.ConsistencyViolation{
			Path:    "domainNames",
			Message: i18n.M(ctx, i18n.K.CommonValidationAtLeastOneRequired),
		})
	}

	for index, domainName := range request.DomainNames {
		if !constants.TLDPattern.MatchString(domainName) {
			violations = append(violations, validation.ConsistencyViolation{
				Path:    fmt.Sprintf("domainNames[%d]", index),
				Message: i18n.M(ctx, i18n.K.CommonValidationInvalidDomainName),
			})
		}
	}

	return violations
}
