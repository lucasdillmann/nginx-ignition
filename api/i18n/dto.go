package i18n

type dictionariesDTO struct {
	DefaultLanguage string          `json:"defaultLanguage"`
	Dictionaries    []dictionaryDTO `json:"dictionaries"`
}
type dictionaryDTO struct {
	Templates map[string]string `json:"templates"`
	Language  string            `json:"language"`
}
