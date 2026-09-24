package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"golang.org/x/text/language"

	"github.com/lucasdillmann/nginx-ignition/internal/business/core/i18n"
)

func Test_i18nMiddleware(t *testing.T) {
	t.Run("uses the first supported language", func(t *testing.T) {
		ctx, commands := newI18nTestContext(t)
		ctx.Request.Header.Set("Accept-Language", "pt-BR, en-US;q=0.8")

		commands.EXPECT().DefaultLanguage().Return(language.AmericanEnglish)
		commands.EXPECT().Supports(language.BrazilianPortuguese).Return(true)

		i18nMiddleware(commands)(ctx)

		assertI18nLanguage(t, ctx, language.BrazilianPortuguese)
	})

	t.Run("selects a later supported language and stops checking", func(t *testing.T) {
		ctx, commands := newI18nTestContext(t)
		ctx.Request.Header.Set("Accept-Language", "fr-FR,de-DE,pt-BR")

		commands.EXPECT().DefaultLanguage().Return(language.AmericanEnglish)
		gomock.InOrder(
			commands.EXPECT().Supports(language.Make("fr-FR")).Return(false),
			commands.EXPECT().Supports(language.Make("de-DE")).Return(false),
			commands.EXPECT().Supports(language.BrazilianPortuguese).Return(true),
		)

		i18nMiddleware(commands)(ctx)

		assertI18nLanguage(t, ctx, language.BrazilianPortuguese)
	})

	t.Run("respects quality value ordering", func(t *testing.T) {
		ctx, commands := newI18nTestContext(t)
		ctx.Request.Header.Set("Accept-Language", "en-US;q=0.1,pt-BR;q=0.9")

		commands.EXPECT().DefaultLanguage().Return(language.AmericanEnglish)
		commands.EXPECT().Supports(language.BrazilianPortuguese).Return(true)

		i18nMiddleware(commands)(ctx)

		assertI18nLanguage(t, ctx, language.BrazilianPortuguese)
	})

	t.Run("falls back when no language is supported", func(t *testing.T) {
		ctx, commands := newI18nTestContext(t)
		ctx.Request.Header.Set("Accept-Language", "fr-FR,de-DE")

		commands.EXPECT().DefaultLanguage().Return(language.AmericanEnglish)
		gomock.InOrder(
			commands.EXPECT().Supports(language.Make("fr-FR")).Return(false),
			commands.EXPECT().Supports(language.Make("de-DE")).Return(false),
		)

		i18nMiddleware(commands)(ctx)

		assertI18nLanguage(t, ctx, language.AmericanEnglish)
	})

	t.Run("falls back when the header is missing", func(t *testing.T) {
		ctx, commands := newI18nTestContext(t)

		commands.EXPECT().DefaultLanguage().Return(language.AmericanEnglish)

		i18nMiddleware(commands)(ctx)

		assertI18nLanguage(t, ctx, language.AmericanEnglish)
	})

	t.Run("falls back when the header is empty", func(t *testing.T) {
		ctx, commands := newI18nTestContext(t)
		ctx.Request.Header.Set("Accept-Language", "")

		commands.EXPECT().DefaultLanguage().Return(language.AmericanEnglish)

		i18nMiddleware(commands)(ctx)

		assertI18nLanguage(t, ctx, language.AmericanEnglish)
	})

	t.Run("falls back when the header is malformed", func(t *testing.T) {
		ctx, commands := newI18nTestContext(t)
		ctx.Request.Header.Set("Accept-Language", "en;q=not-a-number")

		commands.EXPECT().DefaultLanguage().Return(language.AmericanEnglish)

		i18nMiddleware(commands)(ctx)

		assertI18nLanguage(t, ctx, language.AmericanEnglish)
	})

	t.Run("stores the selected language before downstream handlers", func(t *testing.T) {
		ctx, commands := newI18nTestContext(t)
		ctx.Request.Header.Set("Accept-Language", "pt-BR")

		commands.EXPECT().DefaultLanguage().Return(language.AmericanEnglish)
		commands.EXPECT().Supports(language.BrazilianPortuguese).Return(true)

		called := false
		engine := gin.New()
		engine.Use(i18nMiddleware(commands))
		engine.GET("/", func(downstream *gin.Context) {
			called = true
			assertI18nLanguage(t, downstream, language.BrazilianPortuguese)
		})

		recorder := httptest.NewRecorder()
		engine.ServeHTTP(recorder, ctx.Request)

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.True(t, called)
	})
}

func Test_i18nMiddlewareLimits(t *testing.T) {
	t.Run("accepts a header at the byte limit", func(t *testing.T) {
		ctx, commands := newI18nTestContext(t)
		base := strings.Repeat("aa,", maximumLanguageTags) + "aa"
		header := base + strings.Repeat(" ", maximumHeaderBytes-len(base))
		require.Len(t, header, maximumHeaderBytes)
		ctx.Request.Header.Set("Accept-Language", header)

		commands.EXPECT().DefaultLanguage().Return(language.AmericanEnglish)
		expectUnsupportedTags(t, commands, header)

		i18nMiddleware(commands)(ctx)

		assertI18nLanguage(t, ctx, language.AmericanEnglish)
	})

	t.Run("rejects a header above the byte limit", func(t *testing.T) {
		ctx, commands := newI18nTestContext(t)
		base := strings.Repeat("aa,", maximumLanguageTags) + "aa"
		header := base + strings.Repeat(" ", maximumHeaderBytes-len(base)+1)
		require.Len(t, header, maximumHeaderBytes+1)
		ctx.Request.Header.Set("Accept-Language", header)

		commands.EXPECT().DefaultLanguage().Return(language.AmericanEnglish)

		i18nMiddleware(commands)(ctx)

		assertI18nLanguage(t, ctx, language.AmericanEnglish)
	})

	t.Run("accepts the maximum comma count", func(t *testing.T) {
		ctx, commands := newI18nTestContext(t)
		header := strings.Repeat("aa,", maximumLanguageTags) + "aa"
		ctx.Request.Header.Set("Accept-Language", header)

		commands.EXPECT().DefaultLanguage().Return(language.AmericanEnglish)
		expectUnsupportedTags(t, commands, header)

		i18nMiddleware(commands)(ctx)

		assertI18nLanguage(t, ctx, language.AmericanEnglish)
	})

	t.Run("rejects more than the maximum comma count", func(t *testing.T) {
		ctx, commands := newI18nTestContext(t)
		header := strings.Repeat("aa,", maximumLanguageTags+1) + "aa"
		ctx.Request.Header.Set("Accept-Language", header)

		commands.EXPECT().DefaultLanguage().Return(language.AmericanEnglish)

		i18nMiddleware(commands)(ctx)

		assertI18nLanguage(t, ctx, language.AmericanEnglish)
	})

	t.Run("accepts the maximum combined language separator count", func(t *testing.T) {
		ctx, commands := newI18nTestContext(t)
		tags := make([]string, maximumLanguageTags)
		for index := range tags {
			if index%2 == 0 {
				tags[index] = "en-US"
			} else {
				tags[index] = "en_US"
			}
		}
		header := strings.Join(tags, ",")
		ctx.Request.Header.Set("Accept-Language", header)

		commands.EXPECT().DefaultLanguage().Return(language.AmericanEnglish)
		expectUnsupportedTags(t, commands, header)

		i18nMiddleware(commands)(ctx)

		assertI18nLanguage(t, ctx, language.AmericanEnglish)
	})

	t.Run("rejects more than the maximum combined language separator count", func(t *testing.T) {
		ctx, commands := newI18nTestContext(t)
		tags := make([]string, maximumLanguageTags+1)
		for index := range tags {
			if index%2 == 0 {
				tags[index] = "en-US"
			} else {
				tags[index] = "en_US"
			}
		}
		header := strings.Join(tags, ",")
		ctx.Request.Header.Set("Accept-Language", header)

		commands.EXPECT().DefaultLanguage().Return(language.AmericanEnglish)

		i18nMiddleware(commands)(ctx)

		assertI18nLanguage(t, ctx, language.AmericanEnglish)
	})
}

func newI18nTestContext(t *testing.T) (*gin.Context, *i18n.MockedCommands) {
	t.Helper()
	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	commands := i18n.NewMockedCommands(gomock.NewController(t))

	return ctx, commands
}

func assertI18nLanguage(t *testing.T, ctx *gin.Context, expected language.Tag) {
	t.Helper()
	storedLanguage, exists := ctx.Get(i18n.ContextKey)
	assert.Equal(t, expected, ctx.Request.Context().Value(i18n.ContextKey))
	assert.True(t, exists)
	assert.Equal(t, expected, storedLanguage)
}

func expectUnsupportedTags(t *testing.T, commands *i18n.MockedCommands, header string) {
	t.Helper()

	tags, _, err := language.ParseAcceptLanguage(header)
	require.NoError(t, err)
	require.NotEmpty(t, tags)
	for _, tag := range tags {
		commands.EXPECT().Supports(tag).Return(false)
	}
}
