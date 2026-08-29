package settings

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"nginx-ignition/internal/core/settings"
)

type getHandler struct {
	commands settings.Commands
}

func (h getHandler) handle(ctx *gin.Context) {
	data, err := h.commands.Get(ctx.Request.Context())
	if err != nil {
		panic(err)
	}

	ctx.JSON(http.StatusOK, toDTO(data))
}
