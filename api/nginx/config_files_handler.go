package nginx

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"dillmann.com.br/nginx-ignition/core/nginx"
)

type configFilesHandler struct {
	commands *nginx.Commands
}

func (h configFilesHandler) handle(ctx *gin.Context) {
	basePath := normalizePathQuery(ctx, "basePath")
	configPath := normalizePathQuery(ctx, "configPath")
	logPath := normalizePathQuery(ctx, "logPath")

	bytes, err := h.commands.GetConfigFiles(ctx.Request.Context(), basePath, configPath, logPath)
	if err != nil {
		panic(err)
	}

	ctx.Header("Content-Disposition", "attachment; filename=nginx-config.zip")
	ctx.Data(http.StatusOK, "application/zip", bytes)
}

func normalizePathQuery(ctx *gin.Context, name string) string {
	value := ctx.Query(name)
	if value == "" {
		return ""
	}

	return strings.TrimRight(value, "/") + "/"
}
