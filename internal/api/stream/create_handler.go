package stream

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/lucasdillmann/nginx-ignition/internal/api/common/converter"
	"github.com/lucasdillmann/nginx-ignition/internal/core/stream"
)

type createHandler struct {
	commands stream.Commands
}

func (h createHandler) handle(ctx *gin.Context) {
	payload := &streamRequestDTO{}
	if err := ctx.BindJSON(payload); err != nil {
		panic(err)
	}

	domainModel := converter.Wrap(ctx.Request.Context(), toDomain, payload)
	domainModel.ID = uuid.New()

	if err := h.commands.Save(ctx.Request.Context(), domainModel); err != nil {
		panic(err)
	}

	ctx.JSON(
		http.StatusCreated,
		map[string]any{
			"id": domainModel.ID,
		},
	)
}
