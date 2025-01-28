package dynamic_fields

func RemoveSensitiveFields(values *map[string]interface{}, dynamicFields []*DynamicField) {
	for _, field := range dynamicFields {
		if field.Sensitive {
			delete(*values, field.ID)
		}
	}
}
