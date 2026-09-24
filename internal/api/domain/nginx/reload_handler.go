package nginx

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/lucasdillmann/nginx-ignition/internal/business/core/log"
	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/nginx"
)

type reloadHandler struct {
	commands nginx.Commands
}

func (h reloadHandler) handle(ctx *gin.Context) {
	if err := h.commands.Reload(ctx.Request.Context(), false); err != nil {
		log.Warnf("Failed to reload Nginx: %s", err.Error())
		ctx.JSON(http.StatusFailedDependency, gin.H{"message": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}
