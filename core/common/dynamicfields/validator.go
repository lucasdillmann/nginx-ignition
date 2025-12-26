package dynamicfields

import (
	"encoding/base64"
	"net/mail"
	"net/url"
	"strings"

	"dillmann.com.br/nginx-ignition/core/common/ptr"
	"dillmann.com.br/nginx-ignition/core/common/validation"
)

func Validate(
	dynamicFields []DynamicField,
	parameters map[string]any,
) *validation.ConsistencyError {
	var violations []validation.ConsistencyViolation
	for _, field := range dynamicFields {
		value, exists := parameters[field.ID]
		if !field.Required && (!exists || value == "" || value == nil) {
			continue
		}

		conditionSatisfied := areConditionsSatisfied(field, parameters)
		if !exists && field.Required && conditionSatisfied {
			violations = append(violations, validation.ConsistencyViolation{
				Path:    "parameters." + field.ID,
				Message: "A value is required",
			})
		}

		if exists && conditionSatisfied {
			incompatibleMessage := resolveErrorMessage(field, value)
			if incompatibleMessage != nil {
				violations = append(violations, validation.ConsistencyViolation{
					Path:    "parameters." + field.ID,
					Message: *incompatibleMessage,
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
	if field.Conditions == nil || len(field.Conditions) == 0 {
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

func resolveErrorMessage(field DynamicField, value interface{}) *string {
	switch field.Type {
	case EnumType, SingleLineTextType, MultiLineTextType:
		return resolveTextBasedFieldErrorMessage(field, value)

	case FileType:
		if !canDecodeFile(value) {
			return ptr.Of("A file is expected, encoded in a Base64 String")
		}

	case BooleanType:
		if _, ok := value.(bool); !ok {
			return ptr.Of("A boolean value is expected")
		}

	case EmailType:
		if !isAnEmail(value) {
			return ptr.Of("An email is expected")
		}

	case URLType:
		if !isAnUrl(value) {
			return ptr.Of("Not a valid URL")
		}

	default:
		return ptr.Of("Unknown field type")
	}

	return nil
}

func canDecodeFile(value interface{}) bool {
	if value == nil {
		return false
	}

	_, err := base64.StdEncoding.DecodeString(value.(string))
	return err == nil
}

func isAnEmail(value interface{}) bool {
	if value == nil {
		return false
	}

	_, err := mail.ParseAddress(value.(string))
	return err == nil
}

func isAnUrl(value interface{}) bool {
	if value == nil {
		return false
	}

	_, err := url.ParseRequestURI(value.(string))
	return err == nil
}

func resolveTextBasedFieldErrorMessage(field DynamicField, value interface{}) *string {
	castedValue, casted := value.(string)
	if !casted {
		return ptr.Of("A text value is expected")
	}

	if field.Required && strings.TrimSpace(castedValue) == "" {
		return ptr.Of("A not empty text value is required")
	}

	if field.Type == EnumType {
		return resolveEnumFieldErrorMessage(field, castedValue)
	}

	return nil
}

func resolveEnumFieldErrorMessage(field DynamicField, value interface{}) *string {
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
		return ptr.Of("Not a recognized option. Valid values: " + strings.Join(enumOptions, ", "))
	}

	return nil
}
