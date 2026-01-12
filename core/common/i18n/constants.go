package i18n

import (
	"dillmann.com.br/nginx-ignition/core/common/i18n/dict"
)

type LanguageContextKey string

const (
	ContextKey LanguageContextKey = "i18n.language"
)

var K = dict.Keys
