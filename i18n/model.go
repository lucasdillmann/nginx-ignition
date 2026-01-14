package i18n

import "golang.org/x/text/language"

type Dictionary interface {
	Language() language.Tag
	Translate(messageKey string, variables map[string]any) string
	Templates() map[string]string
}
