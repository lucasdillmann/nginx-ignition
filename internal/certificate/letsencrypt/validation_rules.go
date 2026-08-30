package letsencrypt

import (
	"context"
	"fmt"

	"github.com/lucasdillmann/nginx-ignition/internal/core/certificate"
	"github.com/lucasdillmann/nginx-ignition/internal/core/common/dynamicfields"
	"github.com/lucasdillmann/nginx-ignition/internal/core/common/i18n"
	"github.com/lucasdillmann/nginx-ignition/internal/core/common/validation"
)

type validationRules struct {
	dynamicFields []dynamicfields.DynamicField
}

func (r validationRules) DynamicFields() []dynamicfields.DynamicField {
	return r.dynamicFields
}

func (r validationRules) Validate(
	ctx context.Context,
	request *certificate.IssueRequest,
) []validation.ConsistencyViolation {
	output := make([]validation.ConsistencyViolation, 0)

	termsOfServiceAccepted, casted := request.Parameters[termsOfServiceFieldID].(bool)
	if !casted || !termsOfServiceAccepted {
		output = append(output, validation.ConsistencyViolation{
			Path:    fmt.Sprintf("parameters.%s", termsOfServiceFieldID),
			Message: i18n.M(ctx, i18n.K.CertificateLetsencryptTosRequired),
		})
	}

	return output
}
