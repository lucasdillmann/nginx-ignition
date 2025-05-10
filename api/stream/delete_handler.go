package stream

import (
	"dillmann.com.br/nginx-ignition/core/stream"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type deleteHandler struct {
	commands *stream.Commands
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
