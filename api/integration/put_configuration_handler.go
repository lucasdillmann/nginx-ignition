package integration

import (
	"dillmann.com.br/nginx-ignition/core/integration"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type putConfigurationHandler struct {
	command *integration.ConfigureByIdCommand
}

func (h putConfigurationHandler) handle(context *gin.Context) {
	id := context.Param("id")
	if id == "" {
		context.Status(http.StatusNotFound)
		return
	}

	payload := &integrationConfigurationRequest{}
	if err := context.BindJSON(payload); err != nil {
		panic(err)
	}

	if err := validator.New().Struct(payload); err != nil {
		panic(err)
	}

	if err := (*h.command)(id, *payload.Enabled, *payload.Parameters); err != nil {
		panic(err)
	}

	context.Status(http.StatusNoContent)
}
