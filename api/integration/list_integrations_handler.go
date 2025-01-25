package integration

import (
	"dillmann.com.br/nginx-ignition/core/integration"
	"github.com/gin-gonic/gin"
	"net/http"
)

type listIntegrationsHandler struct {
	command *integration.ListCommand
}

func (h listIntegrationsHandler) handle(context *gin.Context) {
	integrations, err := (*h.command)()
	if err != nil {
		panic(err)
	}

	output := make([]*integrationResponse, len(integrations))
	for i, value := range integrations {
		output[i] = toDto(value)
	}

	context.JSON(http.StatusOK, output)
}
