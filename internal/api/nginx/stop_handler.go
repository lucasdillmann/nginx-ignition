package nginx

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"nginx-ignition/internal/core/common/log"
	"nginx-ignition/internal/core/nginx"
)

type stopHandler struct {
	commands nginx.Commands
}

func (h stopHandler) handle(ctx *gin.Context) {
	if err := h.commands.Stop(ctx.Request.Context()); err != nil {
		log.Warnf("Failed to stop Nginx: %s", err.Error())
		ctx.JSON(http.StatusFailedDependency, gin.H{"message": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}
