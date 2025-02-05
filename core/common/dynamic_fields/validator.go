package dynamic_fields

import (
	"dillmann.com.br/nginx-ignition/core/common/validation"
	"encoding/base64"
	"net/mail"
	"net/url"
	"strings"
)

func Validate(
	dynamicFields []*DynamicField,
	parameters map[string]any,
) error {
	var violations []validation.ConsistencyViolation
	for _, field := range dynamicFields {
		value, exists := parameters[field.ID]
		conditionSatisfied := isConditionSatisfied(field, parameters)

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

func isConditionSatisfied(field *DynamicField, parameters map[string]any) bool {
	if field.Condition == nil {
		return true
	}

	expectedValue := field.Condition.Value
	currentValue, exists := parameters[field.Condition.ParentField]
	return exists && expectedValue == currentValue
}

func resolveErrorMessage(field *DynamicField, value interface{}) *string {
	switch field.Type {
	case EnumType, SingleLineTextType, MultiLineTextType:
		return resolveTextBasedFieldErrorMessage(field, value)

	case FileType:
		if !canDecodeFile(value) {
			msg := "A file is expected, encoded in a Base64 String"
			return &msg
		}

	case BooleanType:
		if _, ok := value.(bool); !ok {
			msg := "A boolean value is expected"
			return &msg
		}

	case EmailType:
		if !isAnEmail(value) {
			msg := "An email is expected"
			return &msg
		}

	case URLType:
		if !isAnUrl(value) {
			msg := "A URL is expected"
			return &msg
		}

	default:
		msg := "Unknown field type"
		return &msg
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

func resolveTextBasedFieldErrorMessage(field *DynamicField, value interface{}) *string {
	castedValue, casted := value.(string)
	if !casted {
		msg := "A text value is expected"
		return &msg
	}

	if field.Required && strings.TrimSpace(castedValue) == "" {
		msg := "A not empty text value is required"
		return &msg
	}

	if field.Type == EnumType {
		return resolveEnumFieldErrorMessage(field, castedValue)
	}

	return nil
}

func resolveEnumFieldErrorMessage(field *DynamicField, value interface{}) *string {
	enumOptions := make([]string, len(*field.EnumOptions))
	for i, option := range *field.EnumOptions {
		enumOptions[i] = option.ID
	}

	valid := false
	for _, option := range enumOptions {
		if option == value {
			valid = true
			break
		}
	}

	if !valid {
		msg := "Not a recognized option. Valid values: " + strings.Join(enumOptions, ", ")
		return &msg
	}

	return nil
}
