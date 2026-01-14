package i18n

import (
	"golang.org/x/text/language"

	"dillmann.com.br/nginx-ignition/i18n/dict"
)

type service struct {
	defaultDictionary dict.Dictionary
	dictionaries      []dict.Dictionary
}

func newCommands() Commands {
	dictionaries := []dict.Dictionary{dict.EnUS(), dict.PtBR(), dict.EsES()}
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

func (s *service) resolveDict(lang language.Tag) dict.Dictionary {
	var partialMatch dict.Dictionary

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

func (s *service) GetDictionaries() []dict.Dictionary {
	return s.dictionaries
}
