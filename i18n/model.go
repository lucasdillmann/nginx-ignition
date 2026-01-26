package i18n

import (
	"fmt"
	"os"

	"golang.org/x/text/language"

	"dillmann.com.br/nginx-ignition/core/common/log"
)

type Dictionary interface {
	Language() language.Tag
	Translate(messageKey string, variables map[string]any) string
	Raw() map[string]string
}

type dictionary struct {
	lang     language.Tag
	messages map[string]string
}

func newDictionary(lang language.Tag, messages map[string]string) Dictionary {
	return &dictionary{
		lang:     lang,
		messages: messages,
	}
}

func (d *dictionary) Language() language.Tag {
	return d.lang
}

func (d *dictionary) Raw() map[string]string {
	return d.messages
}

func (d *dictionary) Translate(messageKey string, variables map[string]any) string {
	template, found := d.messages[messageKey]
	if !found {
		log.Warnf(
			"i18n template not found for key %s on language %s. Using the key itself as the message as a fallback.",
			messageKey,
			d.lang.String(),
		)

		return messageKey
	}

	return os.Expand(template, func(varKey string) string {
		if variables == nil || variables[varKey] == nil {
			log.Warnf(
				"i18n variable %s not provided for message %s. Message will be built without replacement as a fallback.",
				varKey,
				messageKey,
			)

			return "${" + varKey + "}"
		}

		return fmt.Sprintf("%v", variables[varKey])
	})
}
