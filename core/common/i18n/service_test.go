package i18n

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/text/language"

	"dillmann.com.br/nginx-ignition/i18n"
)

func Test_service(t *testing.T) {
	t.Run("DefaultLanguage", func(t *testing.T) {
		s := newCommands().(*service)
		assert.Equal(t, language.English, s.DefaultLanguage())
	})

	t.Run("GetDictionaries", func(t *testing.T) {
		s := newCommands().(*service)
		dicts := s.GetDictionaries()
		assert.NotNil(t, dicts)
	})

	t.Run("Supports", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		enUS := i18n.NewMockedDictionary(ctrl)
		ptBR := i18n.NewMockedDictionary(ctrl)

		enUS.EXPECT().Language().Return(language.AmericanEnglish).AnyTimes()
		ptBR.EXPECT().Language().Return(language.BrazilianPortuguese).AnyTimes()

		s := &service{
			dictionaries:      []i18n.Dictionary{enUS, ptBR},
			defaultDictionary: enUS,
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
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		enUS := i18n.NewMockedDictionary(ctrl)
		ptBR := i18n.NewMockedDictionary(ctrl)

		enUS.EXPECT().Language().Return(language.AmericanEnglish).AnyTimes()
		ptBR.EXPECT().Language().Return(language.BrazilianPortuguese).AnyTimes()

		s := &service{
			dictionaries:      []i18n.Dictionary{enUS, ptBR},
			defaultDictionary: enUS,
		}

		t.Run("translates for exact match", func(t *testing.T) {
			enUS.EXPECT().Translate(K.CommonValueMissing, gomock.Nil()).Return("Hello")
			ptBR.EXPECT().Translate(K.CommonValueMissing, gomock.Nil()).Return("Olá")

			assert.Equal(
				t,
				"Hello",
				s.Translate(language.AmericanEnglish, K.CommonValueMissing, nil),
			)
			assert.Equal(
				t,
				"Olá",
				s.Translate(language.BrazilianPortuguese, K.CommonValueMissing, nil),
			)
		})

		t.Run("translates for base language match", func(t *testing.T) {
			ptBR.EXPECT().Translate(K.CommonValueMissing, gomock.Nil()).Return("Olá")

			assert.Equal(
				t,
				"Olá",
				s.Translate(language.EuropeanPortuguese, K.CommonValueMissing, nil),
			)
		})

		t.Run("falls back to default language for unsupported language", func(t *testing.T) {
			enUS.EXPECT().Translate(K.CommonValueMissing, gomock.Nil()).Return("Hello")

			assert.Equal(
				t,
				"Hello",
				s.Translate(language.French, K.CommonValueMissing, nil),
			)
		})

		t.Run("falls back to key if template is missing", func(t *testing.T) {
			enUS.EXPECT().Translate("missing-key", gomock.Nil()).Return("missing-key")

			assert.Equal(
				t,
				"missing-key",
				s.Translate(language.AmericanEnglish, "missing-key", nil),
			)
		})

		t.Run("replaces variables", func(t *testing.T) {
			vars := map[string]any{"name": "Lucas"}
			enUS.EXPECT().Translate(K.CommonBetweenValues, vars).Return("Hello Lucas")

			assert.Equal(
				t,
				"Hello Lucas",
				s.Translate(language.AmericanEnglish, K.CommonBetweenValues, vars),
			)
		})

		t.Run("falls back for missing variables", func(t *testing.T) {
			enUS.EXPECT().
				Translate(K.CommonBetweenValues, gomock.Nil()).
				Return("Hello ${name}")
			enUS.EXPECT().
				Translate(K.CommonBetweenValues, gomock.Any()).
				Return("Hello ${name}")

			assert.Equal(
				t,
				"Hello ${name}",
				s.Translate(language.AmericanEnglish, K.CommonBetweenValues, nil),
			)
			assert.Equal(
				t,
				"Hello ${name}",
				s.Translate(
					language.AmericanEnglish,
					K.CommonBetweenValues,
					make(map[string]any),
				),
			)
		})
	})
}
