package notification

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/notification"
)

type getHandler struct {
	commands notification.Commands
}

func (h getHandler) handle(ctx *gin.Context) {
	userID, ok := currentUserID(ctx)
	if !ok {
		return
	}

	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	data, err := h.commands.GetNotification(ctx.Request.Context(), userID, id)
	if err != nil {
		panic(err)
	}

	if data == nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	ctx.JSON(http.StatusOK, toNotificationDTO(data))
}
