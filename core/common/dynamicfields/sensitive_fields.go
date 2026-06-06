package dynamicfields

func RemoveSensitiveFields(values *map[string]any, dynamicFields []DynamicField) {
	for _, field := range dynamicFields {
		if field.Sensitive {
			delete(*values, field.ID)
		}
	}
}

func MergeSensitiveFields(
	left map[string]any,
	right map[string]any,
	dynamicFields []DynamicField,
) map[string]any {
	if left == nil {
		left = map[string]any{}
	}

	for _, field := range dynamicFields {
		if !field.Sensitive {
			continue
		}

		if _, exists := left[field.ID]; exists {
			continue
		}

		if value, found := right[field.ID]; found {
			left[field.ID] = value
		}
	}

	return left
}
