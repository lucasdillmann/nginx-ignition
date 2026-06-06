package notification

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"dillmann.com.br/nginx-ignition/core/notification"
)

type availableProvidersHandler struct {
	commands notification.Commands
}

func (h availableProvidersHandler) handle(ctx *gin.Context) {
	data, err := h.commands.GetAvailableProviders(ctx.Request.Context())
	if err != nil {
		panic(err)
	}

	payload := make([]availableProviderResponse, len(data))
	for index, item := range data {
		payload[index] = toAvailableProviderDTO(item)
	}

	ctx.JSON(http.StatusOK, payload)
}
