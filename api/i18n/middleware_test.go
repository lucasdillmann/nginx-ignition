package i18n

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/text/language"

	"dillmann.com.br/nginx-ignition/core/i18n"
)

func Test_middleware(t *testing.T) {
	t.Run("middleware", func(t *testing.T) {
		setup := func(t *testing.T) (*gomock.Controller, *i18n.MockedCommands, *gin.Context) {
			ctrl := gomock.NewController(t)
			recorder := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(recorder)
			ctx.Request = httptest.NewRequest("GET", "/", nil)
			commands := i18n.NewMockedCommands(ctrl)
			return ctrl, commands, ctx
		}

		t.Run("uses supported language from Accept-Language header", func(t *testing.T) {
			ctrl, commands, ctx := setup(t)
			defer ctrl.Finish()

			ctx.Request.Header.Set("Accept-Language", "pt-BR, en-US;q=0.8")
			commands.EXPECT().DefaultLanguage().Return(language.AmericanEnglish)
			commands.EXPECT().Supports(language.BrazilianPortuguese).Return(true)

			middleware(commands)(ctx)

			val := ctx.Request.Context().Value(i18n.ContextKey)
			assert.Equal(t, language.BrazilianPortuguese, val)
		})

		t.Run("uses supported base language from Accept-Language header", func(t *testing.T) {
			ctrl, commands, ctx := setup(t)
			defer ctrl.Finish()

			ctx.Request.Header.Set("Accept-Language", "en-GB")
			commands.EXPECT().DefaultLanguage().Return(language.AmericanEnglish)
			commands.EXPECT().Supports(language.BritishEnglish).Return(true)

			middleware(commands)(ctx)

			val := ctx.Request.Context().Value(i18n.ContextKey)
			assert.Equal(t, language.BritishEnglish, val)
		})

		t.Run("falls back to default language when no value matches", func(t *testing.T) {
			ctrl, commands, ctx := setup(t)
			defer ctrl.Finish()

			ctx.Request.Header.Set("Accept-Language", "fr-FR")
			commands.EXPECT().DefaultLanguage().Return(language.AmericanEnglish)
			commands.EXPECT().Supports(language.Make("fr-FR")).Return(false)

			middleware(commands)(ctx)

			val := ctx.Request.Context().Value(i18n.ContextKey)
			assert.Equal(t, language.AmericanEnglish, val)
		})

		t.Run("falls back to default language when header is missing", func(t *testing.T) {
			ctrl, commands, ctx := setup(t)
			defer ctrl.Finish()

			commands.EXPECT().DefaultLanguage().Return(language.AmericanEnglish)

			middleware(commands)(ctx)

			val := ctx.Request.Context().Value(i18n.ContextKey)
			assert.Equal(t, language.AmericanEnglish, val)
		})
	})
}
