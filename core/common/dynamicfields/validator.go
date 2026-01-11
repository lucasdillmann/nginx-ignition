package dynamicfields

import (
	"context"
	"encoding/base64"
	"net/mail"
	"net/url"
	"strings"

	"dillmann.com.br/nginx-ignition/core/common/i18n"
	"dillmann.com.br/nginx-ignition/core/common/validation"
)

func Validate(
	ctx context.Context,
	dynamicFields []DynamicField,
	parameters map[string]any,
) *validation.ConsistencyError {
	violations := make([]validation.ConsistencyViolation, 0)
	for _, field := range dynamicFields {
		value, exists := parameters[field.ID]
		if !field.Required && (!exists || value == "" || value == nil) {
			continue
		}

		conditionSatisfied := areConditionsSatisfied(field, parameters)
		if !exists && field.Required && conditionSatisfied {
			violations = append(violations, validation.ConsistencyViolation{
				Path:    "parameters." + field.ID,
				Message: i18n.M(ctx, "common.validation.value-missing"),
			})
		}

		if exists && conditionSatisfied {
			incompatibleMessage := resolveErrorMessage(ctx, field, value)
			if incompatibleMessage != nil {
				violations = append(violations, validation.ConsistencyViolation{
					Path:    "parameters." + field.ID,
					Message: incompatibleMessage,
				})
			}
		}
	}

	if len(violations) > 0 {
		return validation.NewError(violations)
	}

	return nil
}

func areConditionsSatisfied(field DynamicField, parameters map[string]any) bool {
	if len(field.Conditions) == 0 {
		return true
	}

	for _, condition := range field.Conditions {
		if !isConditionSatisfied(&condition, parameters) {
			return false
		}
	}

	return true
}

func isConditionSatisfied(condition *Condition, parameters map[string]any) bool {
	expectedValue := condition.Value
	currentValue, exists := parameters[condition.ParentField]
	return exists && expectedValue == currentValue
}

func resolveErrorMessage(ctx context.Context, field DynamicField, value any) *i18n.Message {
	switch field.Type {
	case EnumType, SingleLineTextType, MultiLineTextType:
		return resolveTextBasedFieldErrorMessage(ctx, field, value)

	case FileType:
		if !canDecodeFile(value) {
			return i18n.M(ctx, "dynamicfield.validation.invalid-file-encoded-base64")
		}

	case BooleanType:
		if _, ok := value.(bool); !ok {
			return i18n.M(ctx, "dynamicfield.validation.invalid-boolean")
		}

	case EmailType:
		if !isAnEmail(value) {
			return i18n.M(ctx, "dynamicfield.validation.invalid-email")
		}

	case URLType:
		if !isAnURL(value) {
			return i18n.M(ctx, "common.validation.invalid-url")
		}

	default:
		return i18n.M(ctx, "dynamicfield.validation.unknown-field-type")
	}

	return nil
}

func canDecodeFile(value any) bool {
	if value == nil {
		return false
	}

	_, err := base64.StdEncoding.DecodeString(value.(string))
	return err == nil
}

func isAnEmail(value any) bool {
	if value == nil {
		return false
	}

	_, err := mail.ParseAddress(value.(string))
	return err == nil
}

func isAnURL(value any) bool {
	if value == nil {
		return false
	}

	_, err := url.ParseRequestURI(value.(string))
	return err == nil
}

func resolveTextBasedFieldErrorMessage(
	ctx context.Context,
	field DynamicField,
	value any,
) *i18n.Message {
	castedValue, casted := value.(string)
	if !casted {
		return i18n.M(ctx, "dynamicfield.validation.invalid-text")
	}

	if field.Required && strings.TrimSpace(castedValue) == "" {
		return i18n.M(ctx, "common.validation.cannot-be-empty")
	}

	if field.Type == EnumType {
		return resolveEnumFieldErrorMessage(ctx, field, castedValue)
	}

	return nil
}

func resolveEnumFieldErrorMessage(
	ctx context.Context,
	field DynamicField,
	value any,
) *i18n.Message {
	enumOptions := make([]string, len(field.EnumOptions))
	for index, option := range field.EnumOptions {
		enumOptions[index] = option.ID
	}

	valid := false
	for _, option := range enumOptions {
		if option == value {
			valid = true
			break
		}
	}

	if !valid {
		return i18n.M(ctx, "dynamicfield.validation.not-recognized-option").
			V("options", strings.Join(enumOptions, ", "))
	}

	return nil
}
