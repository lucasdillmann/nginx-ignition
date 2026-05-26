package server

import (
	"context"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/text/language"

	"dillmann.com.br/nginx-ignition/core/common/i18n"
)

const maximumLanguageTags = 10

func i18nMiddleware(commands i18n.Commands) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		lang := commands.DefaultLanguage()

		langHeader := ginCtx.GetHeader("Accept-Language")
		if strings.Count(langHeader, "-")+strings.Count(langHeader, "_") > maximumLanguageTags {
			langHeader = ""
		}

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
