package i18n

type dictionaryResponseDTO struct {
	Messages map[string]string `json:"messages"`
	Language string            `json:"languageTag"`
}

type availableLanguagesResponseDTO struct {
	DefaultLanguage string   `json:"default"`
	Available       []string `json:"available"`
}
