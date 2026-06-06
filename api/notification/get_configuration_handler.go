package notification

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/notification"
)

type getConfigurationHandler struct {
	commands notification.Commands
}

func (h getConfigurationHandler) handle(ctx *gin.Context) {
	userID, ok := currentUserID(ctx)
	if !ok {
		return
	}

	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	data, err := h.commands.GetConfiguration(ctx.Request.Context(), userID, id)
	if err != nil {
		panic(err)
	}

	if data == nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	response, err := toConfigurationDTO(data)
	if err != nil {
		panic(err)
	}

	ctx.JSON(http.StatusOK, response)
}
