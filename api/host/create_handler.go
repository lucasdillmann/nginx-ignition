package host

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/api/common/converter"
	"dillmann.com.br/nginx-ignition/core/host"
)

type createHandler struct {
	commands *host.Commands
}

func (h createHandler) handle(ctx *gin.Context) {
	payload := &hostRequestDto{}
	if err := ctx.BindJSON(payload); err != nil {
		panic(err)
	}

	domainModel := converter.Wrap(toDomain, payload)
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
