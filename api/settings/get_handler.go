package settings

import (
	"dillmann.com.br/nginx-ignition/core/settings"
	"github.com/gin-gonic/gin"
	"net/http"
)

type getHandler struct {
	command *settings.GetCommand
}

func (h getHandler) handle(context *gin.Context) {
	data, err := (*h.command)()
	if err != nil {
		panic(err)
	}

	context.JSON(http.StatusOK, toDto(data))
}
