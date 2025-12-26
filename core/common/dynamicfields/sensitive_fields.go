package dynamicfields

func RemoveSensitiveFields(values *map[string]any, dynamicFields []DynamicField) {
	for _, field := range dynamicFields {
		if field.Sensitive {
			delete(*values, field.ID)
		}
	}
}
