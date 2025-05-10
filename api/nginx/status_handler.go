package nginx

import (
	"dillmann.com.br/nginx-ignition/core/nginx"
	"github.com/gin-gonic/gin"
	"net/http"
)

type statusHandler struct {
	commands *nginx.Commands
}

func (h statusHandler) handle(ctx *gin.Context) {
	running := h.commands.GetStatus(ctx.Request.Context())
	ctx.JSON(http.StatusOK, gin.H{"running": running})
}
