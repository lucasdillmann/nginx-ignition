package dynamic_field

import "dillmann.com.br/nginx-ignition/core/common/dynamic_fields"

func ToResponse(fields []*dynamic_fields.DynamicField) []*DynamicFieldResponse {
	response := make([]*DynamicFieldResponse, len(fields))
	for index, field := range fields {
		response[index] = &DynamicFieldResponse{
			ID:           field.ID,
			Priority:     field.Priority,
			Description:  field.Description,
			Required:     field.Required,
			Sensitive:    field.Sensitive,
			Type:         string(field.Type),
			EnumOptions:  toEnumOptions(field.EnumOptions),
			HelpText:     field.HelpText,
			Condition:    toCondition(field.Condition),
			DefaultValue: field.DefaultValue,
		}
	}
	return response
}

func toEnumOptions(options *[]*dynamic_fields.EnumOption) []EnumOption {
	if options == nil {
		return nil
	}

	enumOptions := make([]EnumOption, len(*options))
	for index, option := range *options {
		enumOptions[index] = EnumOption{
			ID:          option.ID,
			Description: option.Description,
		}
	}

	return enumOptions
}

func toCondition(condition *dynamic_fields.Condition) *Condition {
	if condition == nil {
		return nil
	}

	return &Condition{
		ParentField: condition.ParentField,
		Value:       condition.Value,
	}
}
