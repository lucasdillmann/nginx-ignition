package nginx

import (
	"dillmann.com.br/nginx-ignition/core/nginx"
	"github.com/gin-gonic/gin"
	"net/http"
)

type statusHandler struct {
	command *nginx.GetStatusCommand
}

func (h statusHandler) handle(context *gin.Context) {
	running := (*h.command)()
	context.JSON(http.StatusOK, gin.H{"running": running})
}
