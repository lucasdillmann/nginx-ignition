package integration

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/integration"
)

type getOptionHandler struct {
	commands *integration.Commands
}

func (h getOptionHandler) handle(ctx *gin.Context) {
	integrationID := ctx.Param("id")
	if integrationID == "" {
		ctx.Status(http.StatusNotFound)
		return
	}

	integrationUUID, err := uuid.Parse(integrationID)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	optionID := ctx.Param("optionID")
	if optionID == "" {
		ctx.Status(http.StatusNotFound)
		return
	}

	data, err := h.commands.GetOption(ctx.Request.Context(), integrationUUID, optionID)
	if err != nil {
		panic(err)
	}

	if data == nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	ctx.JSON(http.StatusOK, toOptionDto(data))
}
