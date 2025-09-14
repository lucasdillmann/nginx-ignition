package integration

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"dillmann.com.br/nginx-ignition/core/integration"
)

type listIntegrationsHandler struct {
	commands *integration.Commands
}

func (h listIntegrationsHandler) handle(ctx *gin.Context) {
	integrations, err := h.commands.List(ctx.Request.Context())
	if err != nil {
		panic(err)
	}

	output := make([]*integrationResponse, len(integrations))
	for i, value := range integrations {
		output[i] = toDto(value)
	}

	ctx.JSON(http.StatusOK, output)
}
