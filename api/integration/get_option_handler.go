package integration

import (
	"dillmann.com.br/nginx-ignition/core/integration"
	"github.com/gin-gonic/gin"
	"net/http"
)

type getOptionHandler struct {
	commands *integration.Commands
}

func (h getOptionHandler) handle(ctx *gin.Context) {
	integrationId := ctx.Param("id")
	if integrationId == "" {
		ctx.Status(http.StatusNotFound)
		return
	}

	optionId := ctx.Param("optionId")
	if optionId == "" {
		ctx.Status(http.StatusNotFound)
		return
	}

	data, err := h.commands.GetOptionById(ctx.Request.Context(), integrationId, optionId)
	if err != nil {
		panic(err)
	}

	if data == nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	ctx.JSON(http.StatusOK, toOptionDto(data))
}
