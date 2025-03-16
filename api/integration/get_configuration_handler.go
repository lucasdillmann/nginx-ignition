package integration

import (
	"dillmann.com.br/nginx-ignition/core/integration"
	"github.com/gin-gonic/gin"
	"net/http"
)

type getConfigurationHandler struct {
	command *integration.GetByIdCommand
}

func (h getConfigurationHandler) handle(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.Status(http.StatusNotFound)
		return
	}

	data, err := (*h.command)(ctx.Request.Context(), id)
	if err != nil {
		panic(err)
	}

	if data == nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	ctx.JSON(http.StatusOK, toConfigurationDto(data))
}
