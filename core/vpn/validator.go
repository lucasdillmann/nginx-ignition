package vpn

import (
	"context"
	"strings"

	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
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
		v.delegate.Add("enabled", i18n.M(ctx, i18n.K.CoreVpnCannotDisableInUse))
	}

	if strings.TrimSpace(data.Name) == "" {
		v.delegate.Add("name", i18n.M(ctx, i18n.K.CommonValueMissing))
	}

	if strings.TrimSpace(data.Driver) == "" {
		v.delegate.Add("driver", i18n.M(ctx, i18n.K.CommonValueMissing))
	} else if v.driver == nil {
		v.delegate.Add("driver", i18n.M(ctx, i18n.K.CommonInvalidValue))
	}

	params := data.Parameters
	if params == nil {
		params = map[string]any{}
	}

	if v.driver != nil {
		if err := dynamicfields.Validate(
			ctx,
			v.driver.ConfigurationFields(ctx),
			params,
		); err != nil {
			for _, violation := range err.Violations {
				v.delegate.Add(violation.Path, violation.Message)
			}
		}
	}

	return v.delegate.Result()
}
