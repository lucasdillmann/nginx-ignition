package nginx

import (
	"dillmann.com.br/nginx-ignition/core/nginx"
	"github.com/gin-gonic/gin"
	"net/http"
)

type stopHandler struct {
	command *nginx.StopCommand
}

func (h stopHandler) handle(context *gin.Context) {
	if err := (*h.command)(nil); err != nil {
		context.JSON(http.StatusFailedDependency, gin.H{"message": err.Error()})
		return
	}

	context.Status(http.StatusNoContent)
}
