package nginx

import (
	"dillmann.com.br/nginx-ignition/core/common/log"
	"dillmann.com.br/nginx-ignition/core/nginx"
	"github.com/gin-gonic/gin"
	"net/http"
)

type startHandler struct {
	command *nginx.StartCommand
}

func (h startHandler) handle(context *gin.Context) {
	if err := (*h.command)(); err != nil {
		log.Warnf("Failed to start Nginx: %s", err.Error())
		context.JSON(http.StatusFailedDependency, gin.H{"message": err.Error()})
		return
	}

	context.Status(http.StatusNoContent)
}
