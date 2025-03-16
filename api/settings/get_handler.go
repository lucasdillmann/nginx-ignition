package settings

import (
	"dillmann.com.br/nginx-ignition/core/settings"
	"github.com/gin-gonic/gin"
	"net/http"
)

type getHandler struct {
	command *settings.GetCommand
}

func (h getHandler) handle(ctx *gin.Context) {
	data, err := (*h.command)(ctx.Request.Context())
	if err != nil {
		panic(err)
	}

	ctx.JSON(http.StatusOK, toDto(data))
}
