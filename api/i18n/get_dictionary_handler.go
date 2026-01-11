package i18n

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"dillmann.com.br/nginx-ignition/core/i18n"
)

type getDictionaryHandler struct {
	commands i18n.Commands
}

func (h getDictionaryHandler) handle(ctx *gin.Context) {
	dictionaries := h.commands.GetDictionaries()
	ctx.JSON(http.StatusOK, toDTO(dictionaries))
}
