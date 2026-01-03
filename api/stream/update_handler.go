package stream

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/api/common/converter"
	"dillmann.com.br/nginx-ignition/core/stream"
)

type updateHandler struct {
	commands stream.Commands
}

func (h updateHandler) handle(ctx *gin.Context) {
	payload := &streamRequestDTO{}
	if err := ctx.BindJSON(payload); err != nil {
		panic(err)
	}

	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil || id == uuid.Nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	domainModel := converter.Wrap(toDomain, payload)
	domainModel.ID = id

	if err = h.commands.Save(ctx.Request.Context(), domainModel); err != nil {
		panic(err)
	}

	ctx.Status(http.StatusNoContent)
}
