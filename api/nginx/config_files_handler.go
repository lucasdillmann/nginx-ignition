package nginx

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"dillmann.com.br/nginx-ignition/core/nginx"
)

type configFilesHandler struct {
	commands *nginx.Commands
}

func (h configFilesHandler) handle(ctx *gin.Context) {
	bytes, err := h.commands.GetConfigFiles(ctx.Request.Context())
	if err != nil {
		panic(err)
	}

	ctx.Header("Content-Disposition", "attachment; filename=nginx-config-files.zip")
	ctx.Data(http.StatusOK, "application/zip", bytes)
}
