package vpn

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"dillmann.com.br/nginx-ignition/core/vpn"
)

type availableDriversHandler struct {
	commands *vpn.Commands
}

func (h availableDriversHandler) handle(ctx *gin.Context) {
	data, err := h.commands.GetAvailableDrivers(ctx.Request.Context())
	if err != nil || data == nil {
		panic(err)
	}

	payload := make([]vpnDriverResponse, len(data))
	for index, driver := range data {
		payload[index] = toAvailableDriverDTO(&driver)
	}

	ctx.JSON(http.StatusOK, payload)
}
