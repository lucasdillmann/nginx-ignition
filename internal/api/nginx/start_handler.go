package nginx

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"nginx-ignition/internal/core/common/log"
	"nginx-ignition/internal/core/nginx"
)

type startHandler struct {
	commands nginx.Commands
}

func (h startHandler) handle(ctx *gin.Context) {
	if err := h.commands.Start(ctx.Request.Context()); err != nil {
		log.Warnf("Failed to start Nginx: %s", err.Error())
		ctx.JSON(http.StatusFailedDependency, gin.H{"message": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}
