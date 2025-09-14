package host

import (
	"net/http"

	"dillmann.com.br/nginx-ignition/core/host"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type getHandler struct {
	commands *host.Commands
}

func (h getHandler) handle(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	data, err := h.commands.Get(ctx.Request.Context(), id)
	if err != nil {
		panic(err)
	}

	if data == nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	ctx.JSON(http.StatusOK, toDto(data))
}
