package i18n

import (
	"golang.org/x/text/language"
)

type Commands interface {
	Translate(lang language.Tag, messageKey string, variables map[string]any) string
	GetDictionaries() []Dictionary
	DefaultLanguage() language.Tag
	Supports(lang language.Tag) bool
}
