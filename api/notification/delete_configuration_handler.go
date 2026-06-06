package notification

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/notification"
)

type deleteConfigurationHandler struct {
	commands notification.Commands
}

func (h deleteConfigurationHandler) handle(ctx *gin.Context) {
	userID, ok := currentUserID(ctx)
	if !ok {
		return
	}

	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	if err := h.commands.DeleteConfiguration(ctx.Request.Context(), userID, id); err != nil {
		panic(err)
	}

	ctx.Status(http.StatusNoContent)
}
