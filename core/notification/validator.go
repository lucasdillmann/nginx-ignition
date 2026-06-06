package notification

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
	"dillmann.com.br/nginx-ignition/core/common/validation"
)

type validator struct {
	repository Repository
	provider   Provider
	delegate   *validation.ConsistencyValidator
}

func newValidator(repository Repository, provider Provider) *validator {
	return &validator{
		repository: repository,
		provider:   provider,
		delegate:   validation.NewValidator(),
	}
}

func (v *validator) validate(
	ctx context.Context,
	userID uuid.UUID,
	configuration *Configuration,
) error {
	if strings.TrimSpace(configuration.Name) == "" {
		v.delegate.Add("name", i18n.M(ctx, i18n.K.CommonValueMissing))
	} else {
		exists, err := v.repository.ConfigurationExistsByName(
			ctx,
			userID,
			configuration.Name,
			&configuration.ID,
		)
		if err != nil {
			return err
		}

		if exists {
			v.delegate.Add("name", i18n.M(ctx, i18n.K.CoreNotificationDuplicatedName))
		}
	}

	if strings.TrimSpace(configuration.Provider) == "" {
		v.delegate.Add("provider", i18n.M(ctx, i18n.K.CommonValueMissing))
	} else if v.provider == nil {
		v.delegate.Add("provider", i18n.M(ctx, i18n.K.CommonInvalidValue))
	}

	v.validateCategories(ctx, configuration.Categories)

	params := configuration.Parameters
	if params == nil {
		params = map[string]any{}
	}

	if v.provider != nil {
		if err := dynamicfields.Validate(
			ctx,
			v.provider.ConfigurationFields(ctx),
			params,
		); err != nil {
			for _, violation := range err.Violations {
				v.delegate.Add(violation.Path, violation.Message)
			}
		}
	}

	return v.delegate.Result()
}

func (v *validator) validateCategories(ctx context.Context, categories *[]Category) {
	if categories == nil {
		return
	}

	for index, category := range *categories {
		if !isValidCategory(category) {
			v.delegate.Add(
				fmt.Sprintf("categories[%d]", index),
				i18n.M(ctx, i18n.K.CommonInvalidValue),
			)
		}
	}
}
