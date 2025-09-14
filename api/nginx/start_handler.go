package nginx

import (
	"net/http"

	"dillmann.com.br/nginx-ignition/core/common/log"
	"dillmann.com.br/nginx-ignition/core/nginx"
	"github.com/gin-gonic/gin"
)

type startHandler struct {
	commands *nginx.Commands
}

func (h startHandler) handle(ctx *gin.Context) {
	if err := h.commands.Start(ctx.Request.Context()); err != nil {
		log.Warnf("Failed to start Nginx: %s", err.Error())
		ctx.JSON(http.StatusFailedDependency, gin.H{"message": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}
