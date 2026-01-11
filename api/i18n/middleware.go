package i18n

import (
	"context"

	"github.com/gin-gonic/gin"
	"golang.org/x/text/language"

	"dillmann.com.br/nginx-ignition/core/i18n"
)

func middleware(commands i18n.Commands) gin.HandlerFunc {
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

		updatedCtx := context.WithValue(ginCtx.Request.Context(), i18n.ContextKey, lang)
		ginCtx.Request = ginCtx.Request.WithContext(updatedCtx)
		ginCtx.Next()
	}
}
