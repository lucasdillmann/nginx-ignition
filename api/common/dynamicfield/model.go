package dynamicfield

type DynamicFieldResponse struct {
	ID           string       `json:"id"`
	Priority     int          `json:"priority"`
	Description  string       `json:"description"`
	Required     bool         `json:"required"`
	Sensitive    bool         `json:"sensitive"`
	Type         string       `json:"type"`
	EnumOptions  []EnumOption `json:"enumOptions,omitempty"`
	HelpText     *string      `json:"helpText,omitempty"`
	Condition    *Condition   `json:"condition,omitempty"`
	DefaultValue *string      `json:"defaultValue,omitempty"`
}

type EnumOption struct {
	ID          string `json:"id"`
	Description string `json:"description"`
}

type Condition struct {
	ParentField string `json:"parentField"`
	Value       any    `json:"value"`
}
