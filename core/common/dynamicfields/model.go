package dynamicfields

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
	HelpText     *string
	ID           string
	Description  string
	Type         Type
	EnumOptions  []EnumOption
	Conditions   []Condition
	Priority     int
	Required     bool
	Sensitive    bool
}

type EnumOption struct {
	ID          string
	Description string
}

type Condition struct {
	Value       any
	ParentField string
}
