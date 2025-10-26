package integration

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"dillmann.com.br/nginx-ignition/core/integration"
)

type availableDriversHandler struct {
	commands *integration.Commands
}

func (h availableDriversHandler) handle(ctx *gin.Context) {
	data, err := h.commands.GetAvailableDrivers(ctx.Request.Context())
	if err != nil || data == nil {
		panic(err)
	}

	payload := make([]*integrationDriverResponse, len(*data))
	for i, driver := range *data {
		payload[i] = toAvailableDriverDto(driver)
	}

	ctx.JSON(http.StatusOK, payload)
}
