package i18n

import (
	"golang.org/x/text/language"

	"github.com/lucasdillmann/nginx-ignition/internal/i18n"
)

type Commands interface {
	Translate(lang language.Tag, messageKey string, variables map[string]any) string
	GetDictionaries() []i18n.Dictionary
	DefaultLanguage() language.Tag
	Supports(lang language.Tag) bool
}
