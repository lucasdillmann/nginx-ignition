package settings

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/lucasdillmann/nginx-ignition/internal/api/core/converter"

	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/settings"
)

type putHandler struct {
	commands settings.Commands
}

func (h putHandler) handle(ctx *gin.Context) {
	payload := &settingsDTO{}
	if err := ctx.BindJSON(payload); err != nil {
		panic(err)
	}

	domain := converter.Wrap(ctx.Request.Context(), toDomain, payload)
	if err := h.commands.Save(ctx.Request.Context(), domain); err != nil {
		panic(err)
	}

	ctx.Status(http.StatusNoContent)
}
