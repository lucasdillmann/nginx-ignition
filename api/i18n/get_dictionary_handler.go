package i18n

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/text/language"

	core18n "dillmann.com.br/nginx-ignition/core/common/i18n"
	"dillmann.com.br/nginx-ignition/i18n"
)

type getDictionaryHandler struct {
	commands core18n.Commands
}

func (h getDictionaryHandler) handle(ctx *gin.Context) {
	dictionaries := h.commands.GetDictionaries()
	targetLanguageRaw := ctx.Param("language")

	targetLanguage, err := language.Parse(targetLanguageRaw)
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	var dict *i18n.Dictionary
	for _, dictionary := range dictionaries {
		if dictionary.Language() == targetLanguage {
			dict = &dictionary
			break
		}
	}

	if dict == nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	ctx.JSON(http.StatusOK, toDictionaryDTO(*dict))
}
