package dynamicfields

func newDynamicField() *DynamicField {
	return &DynamicField{
		ID:          "field1",
		Description: "A test field",
		Type:        SingleLineTextType,
		Priority:    100,
		Required:    false,
		Sensitive:   false,
	}
}
