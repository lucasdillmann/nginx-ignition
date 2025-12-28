package settings

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"dillmann.com.br/nginx-ignition/api/common/converter"
	"dillmann.com.br/nginx-ignition/core/settings"
)

type putHandler struct {
	commands *settings.Commands
}

func (h putHandler) handle(ctx *gin.Context) {
	payload := &settingsDTO{}
	if err := ctx.BindJSON(payload); err != nil {
		panic(err)
	}

	domain := converter.Wrap(toDomain, payload)
	if err := h.commands.Save(ctx.Request.Context(), &domain); err != nil {
		panic(err)
	}

	ctx.Status(http.StatusNoContent)
}
