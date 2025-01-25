package integration

import (
	"dillmann.com.br/nginx-ignition/core/integration"
	"github.com/gin-gonic/gin"
	"net/http"
)

type getConfigurationHandler struct {
	command *integration.GetByIdCommand
}

func (h getConfigurationHandler) handle(context *gin.Context) {
	id := context.Param("id")
	if id == "" {
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

	context.JSON(http.StatusOK, toConfigurationDto(data))
}
