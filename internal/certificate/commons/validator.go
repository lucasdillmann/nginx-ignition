package commons

import (
	"context"
	"fmt"

	"github.com/lucasdillmann/nginx-ignition/internal/core/certificate"
	"github.com/lucasdillmann/nginx-ignition/internal/core/common/constants"
	"github.com/lucasdillmann/nginx-ignition/internal/core/common/dynamicfields"
	"github.com/lucasdillmann/nginx-ignition/internal/core/common/i18n"
	"github.com/lucasdillmann/nginx-ignition/internal/core/common/validation"
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
			Message: i18n.M(ctx, i18n.K.CommonAtLeastOneRequired),
		})
	}

	for index, domainName := range request.DomainNames {
		if !constants.TLDPattern.MatchString(domainName) {
			violations = append(violations, validation.ConsistencyViolation{
				Path:    fmt.Sprintf("domainNames[%d]", index),
				Message: i18n.M(ctx, i18n.K.CommonInvalidDomainName),
			})
		}
	}

	return violations
}
