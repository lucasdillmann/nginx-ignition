package i18n

type dictionaryDTO struct {
	Templates map[string]string `json:"templates"`
	Language  string            `json:"language"`
}
