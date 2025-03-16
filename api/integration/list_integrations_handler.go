package integration

import (
	"dillmann.com.br/nginx-ignition/core/integration"
	"github.com/gin-gonic/gin"
	"net/http"
)

type listIntegrationsHandler struct {
	command *integration.ListCommand
}

func (h listIntegrationsHandler) handle(ctx *gin.Context) {
	integrations, err := (*h.command)(ctx.Request.Context())
	if err != nil {
		panic(err)
	}

	output := make([]*integrationResponse, len(integrations))
	for i, value := range integrations {
		output[i] = toDto(value)
	}

	ctx.JSON(http.StatusOK, output)
}
