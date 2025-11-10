package stream

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/stream"
)

type createHandler struct {
	commands *stream.Commands
}

func (h createHandler) handle(ctx *gin.Context) {
	payload := &streamRequestDto{}
	if err := ctx.BindJSON(payload); err != nil {
		panic(err)
	}

	domainModel := toDomain(payload)
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
