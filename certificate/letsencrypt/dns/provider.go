package dns

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"

	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
)

type Provider interface {
	ID() string
	Name() string
	DynamicFields() []dynamicfields.DynamicField
	ChallengeProvider(
		ctx context.Context,
		domainNames []string,
		parameters map[string]any,
	) (challenge.Provider, error)
}

func LinkedToProvider(id string, fields []dynamicfields.DynamicField) []dynamicfields.DynamicField {
	output := make([]dynamicfields.DynamicField, 0, len(fields))

	for index, field := range fields {
		field.Priority = index + 2
		if field.Conditions == nil {
			field.Conditions = make([]dynamicfields.Condition, 0, 1)
		}

		field.Conditions = append(field.Conditions, dynamicfields.Condition{
			ParentField: "challengeDnsProvider",
			Value:       id,
		})

		output = append(output, field)
	}

	return output
}
