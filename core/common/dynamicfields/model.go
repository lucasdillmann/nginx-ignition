package dynamicfields

import (
	"dillmann.com.br/nginx-ignition/core/common/i18n"
)

type Type string

const (
	SingleLineTextType Type = "SINGLE_LINE_TEXT"
	MultiLineTextType  Type = "MULTI_LINE_TEXT"
	URLType            Type = "URL"
	EmailType          Type = "EMAIL"
	BooleanType        Type = "BOOLEAN"
	EnumType           Type = "ENUM"
	FileType           Type = "FILE"
)

type DynamicField struct {
	DefaultValue any
	HelpText     *i18n.Message
	ID           string
	Description  *i18n.Message
	Type         Type
	EnumOptions  []EnumOption
	Conditions   []Condition
	Priority     int
	Required     bool
	Sensitive    bool
}

type EnumOption struct {
	Description *i18n.Message
	ID          string
}

type Condition struct {
	Value       any
	ParentField string
}
