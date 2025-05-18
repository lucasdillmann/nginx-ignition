package dynamic_fields

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
	ID           string
	Priority     int
	Description  string
	Required     bool
	Type         Type
	EnumOptions  *[]*EnumOption
	Sensitive    bool
	Condition    *Condition
	HelpText     *string
	DefaultValue *string
}

type EnumOption struct {
	ID          string
	Description string
}

type Condition struct {
	ParentField string
	Value       string
}
