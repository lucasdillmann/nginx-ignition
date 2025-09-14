package integration

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"dillmann.com.br/nginx-ignition/core/integration"
)

type getConfigurationHandler struct {
	commands *integration.Commands
}

func (h getConfigurationHandler) handle(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.Status(http.StatusNotFound)
		return
	}

	data, err := h.commands.GetById(ctx.Request.Context(), id)
	if err != nil {
		panic(err)
	}

	if data == nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	ctx.JSON(http.StatusOK, toConfigurationDto(data))
}
