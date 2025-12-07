package integration

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/api/common/converter"
	"dillmann.com.br/nginx-ignition/core/integration"
)

type putHandler struct {
	commands *integration.Commands
}

func (h putHandler) handle(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.Status(http.StatusNotFound)
		return
	}

	uuidValue, err := uuid.Parse(id)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	payload := &integrationRequest{}
	if err := ctx.BindJSON(payload); err != nil {
		panic(err)
	}

	data := converter.Wrap2(fromDto, uuidValue, payload)

	if err := h.commands.Save(ctx.Request.Context(), data); err != nil {
		panic(err)
	}

	ctx.Status(http.StatusNoContent)
}
