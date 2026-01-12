package cache

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/api/common/converter"
	"dillmann.com.br/nginx-ignition/core/cache"
)

type updateHandler struct {
	commands cache.Commands
}

func (h updateHandler) handle(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	var dto cacheRequestDTO
	if err := ctx.BindJSON(&dto); err != nil {
		panic(err)
	}

	domain := converter.Wrap2(ctx.Request.Context(), toDomain, id, &dto)
	if err := h.commands.Save(ctx, domain); err != nil {
		panic(err)
	}

	ctx.Status(http.StatusNoContent)
}
