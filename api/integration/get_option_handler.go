package integration

import (
	"dillmann.com.br/nginx-ignition/core/integration"
	"github.com/gin-gonic/gin"
	"net/http"
)

type getOptionHandler struct {
	command *integration.GetOptionByIdCommand
}

func (h getOptionHandler) handle(context *gin.Context) {
	integrationId := context.Param("id")
	if integrationId == "" {
		context.Status(http.StatusNotFound)
		return
	}

	optionId := context.Param("optionId")
	if optionId == "" {
		context.Status(http.StatusNotFound)
		return
	}

	data, err := (*h.command)(integrationId, optionId)
	if err != nil {
		panic(err)
	}

	if data == nil {
		context.Status(http.StatusNotFound)
		return
	}

	context.JSON(http.StatusOK, toOptionDto(data))
}
