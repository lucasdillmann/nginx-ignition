package dynamicfield

import "dillmann.com.br/nginx-ignition/core/common/dynamicfields"

func ToResponse(fields []dynamicfields.DynamicField) []Response {
	response := make([]Response, len(fields))
	for index, field := range fields {
		response[index] = Response{
			ID:           field.ID,
			Priority:     field.Priority,
			Description:  field.Description,
			Required:     field.Required,
			Sensitive:    field.Sensitive,
			Type:         string(field.Type),
			EnumOptions:  toEnumOptions(field.EnumOptions),
			HelpText:     field.HelpText,
			Conditions:   toConditions(field.Conditions),
			DefaultValue: field.DefaultValue,
		}
	}

	return response
}

func toEnumOptions(options []dynamicfields.EnumOption) []EnumOption {
	if options == nil {
		return nil
	}

	enumOptions := make([]EnumOption, len(options))
	for index, option := range options {
		enumOptions[index] = EnumOption{
			ID:          option.ID,
			Description: option.Description,
		}
	}

	return enumOptions
}

func toConditions(condition []dynamicfields.Condition) []Condition {
	if len(condition) == 0 {
		return nil
	}

	result := make([]Condition, len(condition))
	for index, item := range condition {
		result[index] = Condition{
			ParentField: item.ParentField,
			Value:       item.Value,
		}
	}

	return result
}
