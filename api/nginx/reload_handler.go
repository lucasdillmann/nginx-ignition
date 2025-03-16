package nginx

import (
	"dillmann.com.br/nginx-ignition/core/common/log"
	"dillmann.com.br/nginx-ignition/core/nginx"
	"github.com/gin-gonic/gin"
	"net/http"
)

type reloadHandler struct {
	command *nginx.ReloadCommand
}

func (h reloadHandler) handle(ctx *gin.Context) {
	if err := (*h.command)(ctx.Request.Context(), false); err != nil {
		log.Warnf("Failed to reload Nginx: %s", err.Error())
		ctx.JSON(http.StatusFailedDependency, gin.H{"message": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}
