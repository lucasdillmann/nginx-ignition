package external

import (
	"context"
	"os"

	"github.com/lucasdillmann/nginx-ignition/internal/core/certificate"
	"github.com/lucasdillmann/nginx-ignition/internal/core/common/dynamicfields"
	"github.com/lucasdillmann/nginx-ignition/internal/core/common/i18n"
	"github.com/lucasdillmann/nginx-ignition/internal/core/common/validation"
)

type validationRules struct {
	dynamicFields []dynamicfields.DynamicField
}

func (v *validationRules) DynamicFields() []dynamicfields.DynamicField {
	return v.dynamicFields
}

func (v *validationRules) Validate(
	ctx context.Context,
	request *certificate.IssueRequest,
) []validation.ConsistencyViolation {
	params := request.Parameters
	violations := make([]validation.ConsistencyViolation, 0)

	pathFields := []string{
		publicKeyPathFieldID,
		privateKeyPathFieldID,
	}

	for _, fieldID := range pathFields {
		rawValue, found := params[fieldID].(string)
		if !found || rawValue == "" {
			continue
		}

		if err := pathReadable(rawValue); err != nil {
			violations = append(violations, validation.ConsistencyViolation{
				Path:    "parameters." + fieldID,
				Message: i18n.M(ctx, i18n.K.CertificateExternalPathNotReadable).V("path", rawValue),
			})
		}
	}

	chainRaw, found := params[chainPathFieldID].(string)
	if found && chainRaw != "" {
		if err := pathReadable(chainRaw); err != nil {
			violations = append(violations, validation.ConsistencyViolation{
				Path:    "parameters." + chainPathFieldID,
				Message: i18n.M(ctx, i18n.K.CertificateExternalPathNotReadable).V("path", chainRaw),
			})
		}
	}

	return violations
}

func pathReadable(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}

	return file.Close()
}
