package i18n

import (
	"golang.org/x/text/language"
)

var enUS = Dictionary{
	Language: language.Make("en-US"),
	Templates: map[string]string{
		"value-missing": "A value is required",
		"test-msg":      "Just a test: ${test-var}",
	},
}
