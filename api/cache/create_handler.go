package cache

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/api/common/converter"
	"dillmann.com.br/nginx-ignition/core/cache"
)

type createHandler struct {
	commands cache.Commands
}

func (h createHandler) handle(ctx *gin.Context) {
	var dto cacheRequestDTO
	if err := ctx.BindJSON(&dto); err != nil {
		panic(err)
	}

	id := uuid.New()
	domain := converter.Wrap2(ctx.Request.Context(), toDomain, id, &dto)
	if err := h.commands.Save(ctx.Request.Context(), domain); err != nil {
		panic(err)
	}

	ctx.JSON(http.StatusCreated, toResponseDTO(domain))
}
