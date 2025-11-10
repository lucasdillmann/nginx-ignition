package vpn

import (
	"context"
	"strings"

	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/validation"
)

type validator struct {
	repository Repository
	driver     Driver
	delegate   *validation.ConsistencyValidator
}

func newValidator(repository Repository, driver Driver) *validator {
	return &validator{
		repository: repository,
		driver:     driver,
		delegate:   validation.NewValidator(),
	}
}

func (v *validator) validate(ctx context.Context, data *VPN) error {
	inUse, err := v.repository.InUseByID(ctx, data.ID)
	if err != nil {
		return err
	}

	if *inUse && !data.Enabled {
		v.delegate.Add("enabled", "VPN is in use by one or more hosts. It cannot be disabled.")
	}

	if strings.TrimSpace(data.Name) == "" {
		v.delegate.Add("name", validation.ValueMissingMessage)
	}

	if strings.TrimSpace(data.Driver) == "" {
		v.delegate.Add("driver", validation.ValueMissingMessage)
	} else if v.driver == nil {
		v.delegate.Add("driver", "Invalid value")
	}

	params := data.Parameters
	if params == nil {
		params = map[string]any{}
	}

	if v.driver != nil {
		if err := dynamicfields.Validate(v.driver.ConfigurationFields(), params); err != nil {
			for _, violation := range err.Violations {
				v.delegate.Add(violation.Path, violation.Message)
			}
		}
	}

	return v.delegate.Result()
}
