package notification

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"dillmann.com.br/nginx-ignition/core/notification"
)

type categoriesHandler struct {
	commands notification.Commands
}

func (h categoriesHandler) handle(ctx *gin.Context) {
	if _, ok := currentUserID(ctx); !ok {
		return
	}

	data, err := h.commands.ListCategories(ctx.Request.Context())
	if err != nil {
		panic(err)
	}

	payload := make([]categoryResponse, len(data))
	for index, item := range data {
		payload[index] = toCategoryDTO(item)
	}

	ctx.JSON(http.StatusOK, payload)
}
