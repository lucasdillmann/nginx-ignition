package integration

import (
	"context"
	"strings"

	"dillmann.com.br/nginx-ignition/core/common/dynamic_fields"
	"dillmann.com.br/nginx-ignition/core/common/validation"
)

type validator struct {
	integration Repository
	driver      Driver
	delegate    *validation.ConsistencyValidator
}

func newValidator(repository Repository, driver Driver) *validator {
	return &validator{
		integration: repository,
		driver:      driver,
		delegate:    validation.NewValidator(),
	}
}

const (
	invalidValue = "Invalid value"
)

func (v *validator) validate(ctx context.Context, data *Integration) error { //nolint:revive,unused
	if strings.TrimSpace(data.Name) == "" {
		v.delegate.Add("name", validation.ValueMissingMessage)
	}

	if strings.TrimSpace(data.Driver) == "" {
		v.delegate.Add("driver", validation.ValueMissingMessage)
	} else if v.driver == nil {
		v.delegate.Add("driver", invalidValue)
	}

	params := data.Parameters
	if params == nil {
		params = map[string]any{}
	}

	if v.driver != nil {
		if err := dynamic_fields.Validate(v.driver.ConfigurationFields(), params); err != nil {
			for _, violation := range err.Violations {
				v.delegate.Add(violation.Path, violation.Message)
			}
		}
	}

	return v.delegate.Result()
}
