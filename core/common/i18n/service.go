package i18n

import (
	"golang.org/x/text/language"

	"dillmann.com.br/nginx-ignition/i18n"
)

type service struct {
	defaultDictionary i18n.Dictionary
	dictionaries      []i18n.Dictionary
}

func newCommands() Commands {
	dictionaries := []i18n.Dictionary{i18n.En_US()}
	return &service{
		dictionaries:      dictionaries,
		defaultDictionary: dictionaries[0],
	}
}

func (s *service) DefaultLanguage() language.Tag {
	return s.defaultDictionary.Language()
}

func (s *service) Supports(lang language.Tag) bool {
	dictionary := s.resolveDict(lang)

	requestedBaseLang, _ := lang.Base()
	resolvedBaseLang, _ := dictionary.Language().Base()

	return requestedBaseLang == resolvedBaseLang
}

func (s *service) Translate(lang language.Tag, messageKey string, variables map[string]any) string {
	dictionary := s.resolveDict(lang)
	return dictionary.Translate(messageKey, variables)
}

func (s *service) resolveDict(lang language.Tag) i18n.Dictionary {
	var partialMatch i18n.Dictionary

	for _, dictionary := range s.dictionaries {
		if dictionary.Language().String() == lang.String() {
			return dictionary
		}

		dictBase, _ := lang.Base()
		requestedBase, _ := dictionary.Language().Base()
		if dictBase.String() == requestedBase.String() {
			partialMatch = dictionary
		}
	}

	if partialMatch != nil {
		return partialMatch
	}

	return s.defaultDictionary
}

func (s *service) GetDictionaries() []i18n.Dictionary {
	return s.dictionaries
}
