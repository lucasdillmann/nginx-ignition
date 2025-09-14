package integration

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"dillmann.com.br/nginx-ignition/core/integration"
)

type putConfigurationHandler struct {
	commands *integration.Commands
}

func (h putConfigurationHandler) handle(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.Status(http.StatusNotFound)
		return
	}

	payload := &integrationConfigurationRequest{}
	if err := ctx.BindJSON(payload); err != nil {
		panic(err)
	}

	if err := validator.New().Struct(payload); err != nil {
		panic(err)
	}

	if err := h.commands.ConfigureById(ctx.Request.Context(), id, *payload.Enabled, *payload.Parameters); err != nil {
		panic(err)
	}

	ctx.Status(http.StatusNoContent)
}
