package host

import (
	"net/http"

	"dillmann.com.br/nginx-ignition/core/host"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type deleteHandler struct {
	commands *host.Commands
}

func (h deleteHandler) handle(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	err = h.commands.Delete(ctx.Request.Context(), id)
	if err != nil {
		panic(err)
	}

	ctx.Status(http.StatusNoContent)
}
