package i18n

import (
	"golang.org/x/text/language"

	"dillmann.com.br/nginx-ignition/i18n/dict"
)

type Commands interface {
	Translate(lang language.Tag, messageKey string, variables map[string]any) string
	GetDictionaries() []dict.Dictionary
	DefaultLanguage() language.Tag
	Supports(lang language.Tag) bool
}
