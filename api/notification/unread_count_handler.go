package notification

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"dillmann.com.br/nginx-ignition/core/notification"
)

type unreadCountHandler struct {
	commands notification.Commands
}

func (h unreadCountHandler) handle(ctx *gin.Context) {
	userID, ok := currentUserID(ctx)
	if !ok {
		return
	}

	count, err := h.commands.UnreadCount(ctx.Request.Context(), userID)
	if err != nil {
		panic(err)
	}

	ctx.JSON(http.StatusOK, unreadCountResponse{Count: count})
}
