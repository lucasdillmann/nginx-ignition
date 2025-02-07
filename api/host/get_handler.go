package host

import (
	"dillmann.com.br/nginx-ignition/core/host"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type getHandler struct {
	command *host.GetCommand
}

func (h getHandler) handle(context *gin.Context) {
	id, err := uuid.Parse(context.Param("id"))
	if err != nil {
		context.Status(http.StatusNotFound)
		return
	}

	data, err := (*h.command)(id)
	if err != nil {
		panic(err)
	}

	if data == nil {
		context.Status(http.StatusNotFound)
		return
	}

	context.JSON(http.StatusOK, toDto(data))
}
