package i18n

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
)

func Test_service(t *testing.T) {
	t.Run("DefaultLanguage", func(t *testing.T) {
		s := newCommands().(*service)
		assert.Equal(t, language.AmericanEnglish, s.DefaultLanguage())
	})

	t.Run("GetDictionaries", func(t *testing.T) {
		s := newCommands().(*service)
		dicts := s.GetDictionaries()
		assert.NotNil(t, dicts)
	})

	t.Run("Supports", func(t *testing.T) {
		s := &service{
			languages: []Dictionary{
				{Language: language.AmericanEnglish},
				{Language: language.BrazilianPortuguese},
			},
			defaultLanguage: Dictionary{Language: language.AmericanEnglish},
			cache:           make(map[language.Tag]*Dictionary),
		}

		t.Run("returns true for exact match", func(t *testing.T) {
			assert.True(t, s.Supports(language.AmericanEnglish))
			assert.True(t, s.Supports(language.BrazilianPortuguese))
		})

		t.Run("returns true for base language match", func(t *testing.T) {
			assert.True(t, s.Supports(language.BritishEnglish))
			assert.True(t, s.Supports(language.Portuguese))
		})

		t.Run("returns false for no match", func(t *testing.T) {
			assert.False(t, s.Supports(language.French))
		})
	})

	t.Run("Translate", func(t *testing.T) {
		s := &service{
			languages: []Dictionary{
				{
					Language: language.AmericanEnglish,
					Templates: map[string]string{
						"msg":  "Hello",
						"vars": "Hello ${name}",
					},
				},
				{
					Language: language.BrazilianPortuguese,
					Templates: map[string]string{
						"msg": "Olá",
					},
				},
			},
			defaultLanguage: Dictionary{
				Language: language.AmericanEnglish,
				Templates: map[string]string{
					"msg":  "Hello",
					"vars": "Hello ${name}",
				},
			},
			cache: make(map[language.Tag]*Dictionary),
		}

		t.Run("translates for exact match", func(t *testing.T) {
			assert.Equal(t, "Hello", s.Translate(language.AmericanEnglish, "msg", nil))
			assert.Equal(t, "Olá", s.Translate(language.BrazilianPortuguese, "msg", nil))
		})

		t.Run("translates for base language match", func(t *testing.T) {
			assert.Equal(t, "Hello", s.Translate(language.BritishEnglish, "msg", nil))
		})

		t.Run("falls back to default language for unsupported language", func(t *testing.T) {
			assert.Equal(t, "Hello", s.Translate(language.French, "msg", nil))
		})

		t.Run("falls back to key if template is missing", func(t *testing.T) {
			assert.Equal(
				t,
				"missing-key",
				s.Translate(language.AmericanEnglish, "missing-key", nil),
			)
		})

		t.Run("replaces variables", func(t *testing.T) {
			vars := map[string]any{"name": "Lucas"}
			assert.Equal(t, "Hello Lucas", s.Translate(language.AmericanEnglish, "vars", vars))
		})

		t.Run("falls back for missing variables", func(t *testing.T) {
			assert.Equal(t, "Hello ${name}", s.Translate(language.AmericanEnglish, "vars", nil))
			assert.Equal(
				t,
				"Hello ${name}",
				s.Translate(language.AmericanEnglish, "vars", map[string]any{}),
			)
		})
	})
}
