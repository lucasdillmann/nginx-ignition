package dns

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"

	"dillmann.com.br/nginx-ignition/core/common/dynamic_fields"
)

type Provider interface {
	ID() string
	Name() string
	DynamicFields() []*dynamic_fields.DynamicField
	ChallengeProvider(
		ctx context.Context,
		domainNames []string,
		parameters map[string]any,
	) (challenge.Provider, error)
}

func LinkedToProvider(id string, fields []dynamic_fields.DynamicField) []*dynamic_fields.DynamicField {
	output := make([]*dynamic_fields.DynamicField, 0, len(fields))

	for index, field := range fields {
		field.Priority = index + 2
		if field.Condition == nil {
			field.Condition = &dynamic_fields.Condition{
				ParentField: "challengeDnsProvider",
				Value:       id,
			}
		}

		output = append(output, &field)
	}

	return output
}
