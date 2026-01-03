package integration

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/integration"
)

type deleteHandler struct {
	commands integration.Commands
}

func (h deleteHandler) handle(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.Status(http.StatusNotFound)
		return
	}

	uuidValue, err := uuid.Parse(id)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	err = h.commands.Delete(ctx.Request.Context(), uuidValue)
	if err != nil {
		panic(err)
	}

	ctx.Status(http.StatusNoContent)
}
