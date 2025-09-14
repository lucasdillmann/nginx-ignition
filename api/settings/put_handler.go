package settings

import (
	"net/http"

	"dillmann.com.br/nginx-ignition/core/settings"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type putHandler struct {
	commands *settings.Commands
}

func (h putHandler) handle(ctx *gin.Context) {
	payload := &settingsDto{}
	if err := ctx.BindJSON(payload); err != nil {
		panic(err)
	}

	if err := validator.New().Struct(payload); err != nil {
		panic(err)
	}

	domain := toDomain(payload)
	if err := h.commands.Save(ctx.Request.Context(), domain); err != nil {
		panic(err)
	}

	ctx.Status(http.StatusNoContent)
}
