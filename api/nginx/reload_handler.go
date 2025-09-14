package nginx

import (
	"net/http"

	"dillmann.com.br/nginx-ignition/core/common/log"
	"dillmann.com.br/nginx-ignition/core/nginx"
	"github.com/gin-gonic/gin"
)

type reloadHandler struct {
	commands *nginx.Commands
}

func (h reloadHandler) handle(ctx *gin.Context) {
	if err := h.commands.Reload(ctx.Request.Context(), false); err != nil {
		log.Warnf("Failed to reload Nginx: %s", err.Error())
		ctx.JSON(http.StatusFailedDependency, gin.H{"message": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}
