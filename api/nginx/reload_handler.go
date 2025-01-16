package nginx

import (
	"dillmann.com.br/nginx-ignition/core/nginx"
	"github.com/gin-gonic/gin"
	"net/http"
)

type reloadHandler struct {
	command *nginx.ReloadCommand
}

func (h reloadHandler) handle(context *gin.Context) {
	if err := (*h.command)(false); err != nil {
		context.JSON(http.StatusFailedDependency, gin.H{"message": err.Error()})
		return
	}

	context.Status(http.StatusNoContent)
}
