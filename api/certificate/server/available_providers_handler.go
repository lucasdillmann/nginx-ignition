package server

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"dillmann.com.br/nginx-ignition/core/certificate/server"
)

type availableProvidersHandler struct {
	commands *server.Commands
}

func (h availableProvidersHandler) handle(ctx *gin.Context) {
	availableProviders, err := h.commands.AvailableProviders(ctx.Request.Context())
	if err != nil {
		panic(err)
	}

	ctx.JSON(http.StatusOK, toAvailableProviderResponse(availableProviders))
}
