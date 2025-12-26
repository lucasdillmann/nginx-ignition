package cache

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/cache"
)

type deleteHandler struct {
	commands *cache.Commands
}

func (h deleteHandler) handle(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	if err := h.commands.Delete(ctx.Request.Context(), id); err != nil {
		panic(err)
	}

	ctx.Status(http.StatusNoContent)
}
