package i18n

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/text/language"

	"dillmann.com.br/nginx-ignition/core/common/container"
)

func Test_Message(t *testing.T) {
	t.Run("String", func(t *testing.T) {
		t.Run("returns translated string when language is in context", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			container.Init(context.Background())
			commands := NewMockedCommands(ctrl)
			container.Singleton[Commands](commands)

			lang := language.BrazilianPortuguese
			ctx := context.WithValue(context.Background(), ContextKey, lang)
			key := "test-key"
			variables := map[string]any{"var": "val"}
			message := Message{ctx: ctx, Key: key, Variables: variables}

			expected := "translated string"
			commands.EXPECT().Translate(lang, key, variables).Return(expected)

			result := message.String()
			assert.Equal(t, expected, result)
		})

		t.Run(
			"returns translated string with default language when language is missing in context",
			func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				container.Init(context.Background())
				commands := NewMockedCommands(ctrl)
				container.Singleton[Commands](commands)

				ctx := context.Background()
				key := "test-key"
				variables := map[string]any{"var": "val"}
				message := Message{ctx: ctx, Key: key, Variables: variables}

				defaultLang := language.AmericanEnglish
				commands.EXPECT().DefaultLanguage().Return(defaultLang)
				commands.EXPECT().
					Translate(defaultLang, key, variables).
					Return("default translated")

				result := message.String()
				assert.Equal(t, "default translated", result)
			},
		)
	})

	t.Run("MarshalJSON", func(t *testing.T) {
		t.Run("marshals translated string", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			container.Init(context.Background())
			commands := NewMockedCommands(ctrl)
			container.Singleton[Commands](commands)

			lang := language.AmericanEnglish
			ctx := context.WithValue(context.Background(), ContextKey, lang)
			message := Message{ctx: ctx, Key: "key"}

			commands.EXPECT().Translate(lang, "key", gomock.Any()).Return("translated")

			bytes, err := json.Marshal(message)
			assert.NoError(t, err)
			assert.Equal(t, `"translated"`, string(bytes))
		})
	})
}
