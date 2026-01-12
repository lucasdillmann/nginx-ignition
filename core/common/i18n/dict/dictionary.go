package dict

import (
	"fmt"
	"os"
	"reflect"

	"golang.org/x/text/language"

	"dillmann.com.br/nginx-ignition/core/common/log"
)

type dictionary struct {
	templates map[string]string
	lang      language.Tag
	messages  Messages
}

func newDictionary(lang language.Tag, messages Messages) Dictionary {
	dict := &dictionary{
		lang:      lang,
		messages:  messages,
		templates: make(map[string]string),
	}
	dict.initTemplates()

	return dict
}

func (d *dictionary) Language() language.Tag {
	return d.lang
}

func (d *dictionary) Templates() map[string]string {
	return d.templates
}

func (d *dictionary) initTemplates() {
	keysType := reflect.TypeOf(d.messages)
	keysValue := reflect.ValueOf(d.messages)
	kType := reflect.TypeOf(Keys)

	for index := 0; index < keysType.NumField(); index++ {
		field := keysType.Field(index)
		value := keysValue.Field(index).String()
		_, found := kType.FieldByName(field.Name)

		if found {
			keyName := reflect.ValueOf(Keys).FieldByName(field.Name).String()
			d.templates[keyName] = value
		}
	}
}

func (d *dictionary) Translate(messageKey string, variables map[string]any) string {
	template, found := d.templates[messageKey]
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
