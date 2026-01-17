package dynamicfield

import (
	"dillmann.com.br/nginx-ignition/core/common/i18n"
)

type Response struct {
	DefaultValue any           `json:"defaultValue,omitempty"`
	HelpText     *i18n.Message `json:"helpText,omitempty"`
	ID           string        `json:"id"`
	Description  *i18n.Message `json:"description"`
	Type         string        `json:"type"`
	EnumOptions  []EnumOption  `json:"enumOptions,omitempty"`
	Conditions   []Condition   `json:"conditions,omitempty"`
	Priority     int           `json:"priority"`
	Required     bool          `json:"required"`
	Sensitive    bool          `json:"sensitive"`
}

type EnumOption struct {
	Description *i18n.Message `json:"description"`
	ID          string        `json:"id"`
}

type Condition struct {
	Value       any    `json:"value"`
	ParentField string `json:"parentField"`
}
