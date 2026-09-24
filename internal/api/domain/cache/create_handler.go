package cache

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/lucasdillmann/nginx-ignition/internal/api/core/converter"

	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/cache"
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
