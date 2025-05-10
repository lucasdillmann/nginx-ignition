package settings

import (
	"dillmann.com.br/nginx-ignition/core/settings"
	"github.com/gin-gonic/gin"
	"net/http"
)

type getHandler struct {
	commands *settings.Commands
}

func (h getHandler) handle(ctx *gin.Context) {
	data, err := h.commands.Get(ctx.Request.Context())
	if err != nil {
		panic(err)
	}

	ctx.JSON(http.StatusOK, toDto(data))
}
