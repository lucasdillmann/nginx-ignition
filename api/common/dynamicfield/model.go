package dynamicfield

type Response struct {
	DefaultValue any          `json:"defaultValue,omitempty"`
	HelpText     *string      `json:"helpText,omitempty"`
	ID           string       `json:"id"`
	Description  string       `json:"description"`
	Type         string       `json:"type"`
	EnumOptions  []EnumOption `json:"enumOptions,omitempty"`
	Conditions   []Condition  `json:"conditions,omitempty"`
	Priority     int          `json:"priority"`
	Required     bool         `json:"required"`
	Sensitive    bool         `json:"sensitive"`
}

type EnumOption struct {
	ID          string `json:"id"`
	Description string `json:"description"`
}

type Condition struct {
	Value       any    `json:"value"`
	ParentField string `json:"parentField"`
}
