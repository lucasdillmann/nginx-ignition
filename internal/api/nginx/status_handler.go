package nginx

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/lucasdillmann/nginx-ignition/internal/core/nginx"
)

type statusHandler struct {
	commands nginx.Commands
}

func (h statusHandler) handle(ctx *gin.Context) {
	status := h.commands.GetStatus(ctx.Request.Context())
	ctx.JSON(http.StatusOK, gin.H{
		"running":       status.Running,
		"uptimeSeconds": status.UptimeSeconds,
	})
}
