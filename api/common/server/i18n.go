package server

import (
	"context"

	"github.com/gin-gonic/gin"
	"golang.org/x/text/language"

	"dillmann.com.br/nginx-ignition/i18n"
)

func i18nMiddleware(commands i18n.Commands) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		lang := commands.DefaultLanguage()

		langHeader := ginCtx.GetHeader("Accept-Language")
		tags, _, err := language.ParseAcceptLanguage(langHeader)
		if err == nil && len(tags) > 0 {
			for _, tag := range tags {
				if commands.Supports(tag) {
					lang = tag
					break
				}
			}
		}

		//nolint:staticcheck
		updatedCtx := context.WithValue(ginCtx.Request.Context(), i18n.ContextKey, lang)
		ginCtx.Request = ginCtx.Request.WithContext(updatedCtx)
		ginCtx.Set(i18n.ContextKey, lang)
		ginCtx.Next()
	}
}
