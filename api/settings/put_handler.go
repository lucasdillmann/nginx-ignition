package settings

import (
	"dillmann.com.br/nginx-ignition/core/settings"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type putHandler struct {
	command *settings.SaveCommand
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
	if err := (*h.command)(ctx.Request.Context(), domain); err != nil {
		panic(err)
	}

	ctx.Status(http.StatusNoContent)
}
