package i18n

type dictionariesDTO struct {
	DefaultLanguage string          `json:"defaultLanguage"`
	Dictionaries    []dictionaryDTO `json:"dictionaries"`
}
type dictionaryDTO struct {
	Messages map[string]string `json:"messages"`
	Language string            `json:"language"`
}
