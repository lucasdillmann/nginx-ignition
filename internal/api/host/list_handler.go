package host

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/lucasdillmann/nginx-ignition/internal/api/common/pagination"
	"github.com/lucasdillmann/nginx-ignition/internal/core/common/log"
	"github.com/lucasdillmann/nginx-ignition/internal/core/host"
	"github.com/lucasdillmann/nginx-ignition/internal/core/settings"
)

type listHandler struct {
	settingsCommands settings.Commands
	hostCommands     host.Commands
}

func (h listHandler) handle(ctx *gin.Context) {
	pageSize, pageNumber, searchTerms, err := pagination.ExtractPaginationParameters(ctx)
	if err != nil {
		panic(err)
	}

	page, err := h.hostCommands.List(ctx.Request.Context(), pageSize, pageNumber, searchTerms)
	if err != nil {
		panic(err)
	}

	globalSettings, err := h.settingsCommands.Get(ctx.Request.Context())
	if err != nil {
		log.Warnf(
			"Unable to get global settings (%v). Proceeding without the global bindings filled for now.",
			err,
		)
		globalSettings = nil
	}

	ctx.JSON(
		http.StatusOK,
		pagination.Convert(page, func(item *host.Host) *hostResponseDTO {
			return toDTO(item, globalSettings)
		}),
	)
}
