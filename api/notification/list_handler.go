package notification

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"dillmann.com.br/nginx-ignition/api/common/pagination"
	"dillmann.com.br/nginx-ignition/core/notification"
)

type listHandler struct {
	commands notification.Commands
}

func (h listHandler) handle(ctx *gin.Context) {
	userID, ok := currentUserID(ctx)
	if !ok {
		return
	}

	pageSize, pageNumber, searchTerms, err := pagination.ExtractPaginationParameters(ctx)
	if err != nil {
		panic(err)
	}

	page, err := h.commands.ListNotifications(
		ctx.Request.Context(),
		userID,
		pageSize,
		pageNumber,
		searchTerms,
	)
	if err != nil {
		panic(err)
	}

	ctx.JSON(http.StatusOK, pagination.Convert(page, toNotificationDTO))
}
