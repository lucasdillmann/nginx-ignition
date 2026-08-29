package i18n

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/lucasdillmann/nginx-ignition/internal/core/common/i18n"
)

type getAvailableLanguagesHandler struct {
	commands i18n.Commands
}

func (h getAvailableLanguagesHandler) handle(ctx *gin.Context) {
	dictionaries := h.commands.GetDictionaries()
	defaultLanguage := h.commands.DefaultLanguage()

	ctx.JSON(http.StatusOK, toAvailableLanguagesDTO(dictionaries, defaultLanguage))
}
