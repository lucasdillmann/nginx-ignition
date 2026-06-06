package notification

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"dillmann.com.br/nginx-ignition/core/notification"
)

type listConfigurationsHandler struct {
	commands notification.Commands
}

func (h listConfigurationsHandler) handle(ctx *gin.Context) {
	userID, ok := currentUserID(ctx)
	if !ok {
		return
	}

	data, err := h.commands.ListConfigurations(ctx.Request.Context(), userID)
	if err != nil {
		panic(err)
	}

	payload := make([]configurationResponse, len(data))
	for index, item := range data {
		response, err := toConfigurationDTO(&item)
		if err != nil {
			panic(err)
		}
		payload[index] = response
	}

	ctx.JSON(http.StatusOK, payload)
}
