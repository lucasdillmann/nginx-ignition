package host

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/common/log"
	"dillmann.com.br/nginx-ignition/core/host"
	"dillmann.com.br/nginx-ignition/core/settings"
)

type getHandler struct {
	settingsCommands settings.Commands
	hostCommands     host.Commands
}

func (h getHandler) handle(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	data, err := h.hostCommands.Get(ctx.Request.Context(), id)
	if err != nil {
		panic(err)
	}

	if data == nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	globalSettings, err := h.settingsCommands.Get(ctx.Request.Context())
	if err != nil {
		log.Warnf(
			"Unable to get global settings (%v). Proceeding without the global bindings filled for now.",
			err,
		)
		globalSettings = nil
	}

	ctx.JSON(http.StatusOK, toDTO(data, globalSettings))
}
