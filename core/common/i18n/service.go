package i18n

import (
	"fmt"
	"os"

	"golang.org/x/text/language"

	"dillmann.com.br/nginx-ignition/core/common/log"
)

type service struct {
	defaultLanguage Dictionary
	cache           map[language.Tag]*Dictionary
	languages       []Dictionary
}

func newCommands() Commands {
	return &service{
		cache:           make(map[language.Tag]*Dictionary),
		languages:       []Dictionary{enUS},
		defaultLanguage: enUS,
	}
}

func (s *service) DefaultLanguage() language.Tag {
	return s.defaultLanguage.Language
}

func (s *service) Supports(lang language.Tag) bool {
	dict := s.resolveDict(lang)

	requestedBaseLang, _ := lang.Base()
	resolvedBaseLang, _ := dict.Language.Base()

	return requestedBaseLang == resolvedBaseLang
}

func (s *service) Translate(lang language.Tag, messageKey string, variables map[string]any) string {
	dict := s.resolveDict(lang)
	template, found := dict.Templates[messageKey]
	if !found {
		log.Warnf(
			"i18n template not found for key %s on language %s. Using the key itself as the message as a fallback.",
			messageKey,
			lang.String(),
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

func (s *service) resolveDict(lang language.Tag) *Dictionary {
	if dict, found := s.cache[lang]; found {
		return dict
	}

	var dict *Dictionary
	for _, item := range s.languages {
		if item.Language.String() == lang.String() {
			dict = &item
			break
		}
	}

	if dict == nil {
		for _, item := range s.languages {
			itemBaseLang, _ := item.Language.Base()
			wantedBaseLang, _ := lang.Base()

			if itemBaseLang == wantedBaseLang {
				dict = &item
				break
			}
		}
	}

	if dict == nil {
		dict = &s.defaultLanguage
	}

	s.cache[lang] = dict
	return dict
}

func (s *service) GetDictionaries() []Dictionary {
	return s.languages
}
