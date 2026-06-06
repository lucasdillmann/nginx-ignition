//nolint:staticcheck,revive
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

func TestMain(m *testing.M) {
	m.Run()
	container.Shutdown()
}

func Test_Message(t *testing.T) {
	t.Run("String", func(t *testing.T) {
		t.Run("returns translated string when language is in context", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			container.Init(t.Context())
			commands := NewMockedCommands(ctrl)
			container.Singleton[Commands](commands)

			lang := language.BrazilianPortuguese
			ctx := context.WithValue(t.Context(), ContextKey, lang)
			key := "test-key"
			variables := map[string]any{"var": "val"}
			message := Message{
				ctx:             ctx,
				DetachedMessage: DetachedMessage{Key: key, Variables: variables},
			}

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

				container.Init(t.Context())
				commands := NewMockedCommands(ctrl)
				container.Singleton[Commands](commands)

				key := "test-key"
				variables := map[string]any{"var": "val"}
				message := Message{
					ctx:             t.Context(),
					DetachedMessage: DetachedMessage{Key: key, Variables: variables},
				}

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

			container.Init(t.Context())
			commands := NewMockedCommands(ctrl)
			container.Singleton[Commands](commands)

			lang := language.AmericanEnglish
			ctx := context.WithValue(t.Context(), ContextKey, lang)
			message := Message{
				ctx:             ctx,
				DetachedMessage: DetachedMessage{Key: "key"},
			}

			commands.EXPECT().Translate(lang, "key", gomock.Any()).Return("translated")

			bytes, err := json.Marshal(message)
			assert.NoError(t, err)
			assert.Equal(t, `"translated"`, string(bytes))
		})
	})

	t.Run("Detach", func(t *testing.T) {
		t.Run("copies key and variables without context", func(t *testing.T) {
			message := M(t.Context(), "test-key").V("var", "val")

			detached := message.Detach()

			assert.Equal(t, message.Key, detached.Key)
			assert.Equal(t, message.Variables, detached.Variables)

			message.Variables["var"] = "mutated"
			assert.Equal(t, "val", detached.Variables["var"])
		})

		t.Run("detached map is independent from message variables", func(t *testing.T) {
			message := M(t.Context(), "key").V("var", "val")
			detached := message.Detach()

			detached.Variables["var"] = "changed"

			assert.Equal(t, "val", message.Variables["var"])
			assert.Equal(t, "changed", detached.Variables["var"])

			message.Variables["other"] = "added"

			assert.NotContains(t, detached.Variables, "other")
		})

		t.Run("embedded message still translates via String with context", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			container.Init(t.Context())
			commands := NewMockedCommands(ctrl)
			container.Singleton[Commands](commands)

			lang := language.BrazilianPortuguese
			ctx := context.WithValue(t.Context(), ContextKey, lang)
			key := "test-key"
			variables := map[string]any{"var": "val"}
			message := M(ctx, key).V("var", "val")

			commands.EXPECT().Translate(lang, key, variables).Return("translated string")

			assert.Equal(t, "translated string", message.String())
		})
	})
}
