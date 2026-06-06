package notification

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"dillmann.com.br/nginx-ignition/core/notification"
)

type markAllAsReadHandler struct {
	commands notification.Commands
}

func (h markAllAsReadHandler) handle(ctx *gin.Context) {
	userID, ok := currentUserID(ctx)
	if !ok {
		return
	}

	if err := h.commands.MarkAllAsRead(ctx.Request.Context(), userID); err != nil {
		panic(err)
	}

	ctx.Status(http.StatusNoContent)
}
