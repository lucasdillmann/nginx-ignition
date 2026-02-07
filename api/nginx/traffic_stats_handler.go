package nginx

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"dillmann.com.br/nginx-ignition/core/nginx"
)

type trafficStatsHandler struct {
	commands nginx.Commands
}

func (h trafficStatsHandler) handle(ctx *gin.Context) {
	stats, err := h.commands.GetTrafficStats(ctx.Request.Context())
	if err != nil {
		panic(err)
	}

	response := toTrafficStatsResponseDTO(stats)
	ctx.JSON(http.StatusOK, response)
}
