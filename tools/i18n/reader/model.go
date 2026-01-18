package reader

type PropertiesFile struct {
	LanguageTag           string
	NormalizedLanguageTag string
	Messages              []Message
}

type Message struct {
	PropertiesKey string
	CamelCaseKey  string
	SnakeCaseKey  string
	Value         string
}
