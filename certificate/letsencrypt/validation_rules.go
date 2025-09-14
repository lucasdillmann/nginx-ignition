package letsencrypt

import (
	"fmt"

	"dillmann.com.br/nginx-ignition/core/certificate"
	"dillmann.com.br/nginx-ignition/core/common/dynamic_fields"
	"dillmann.com.br/nginx-ignition/core/common/validation"
)

type validationRules struct {
	dynamicFields []*dynamic_fields.DynamicField
}

func (r validationRules) DynamicFields() []*dynamic_fields.DynamicField {
	return r.dynamicFields
}

func (r validationRules) Validate(request *certificate.IssueRequest) []validation.ConsistencyViolation {
	output := make([]validation.ConsistencyViolation, 0)

	termsOfServiceAccepted, casted := request.Parameters[termsOfService.ID].(bool)
	if !casted || !termsOfServiceAccepted {
		output = append(output, validation.ConsistencyViolation{
			Path:    fmt.Sprintf("parameters.%s", termsOfService.ID),
			Message: "You must accept the Let's Encrypt's terms of service to be able to use its certificates",
		})
	}

	return output
}
