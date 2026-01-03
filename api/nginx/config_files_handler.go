package nginx

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"dillmann.com.br/nginx-ignition/core/nginx"
)

type configFilesHandler struct {
	commands nginx.Commands
}

func (h configFilesHandler) handle(ctx *gin.Context) {
	bytes, err := h.commands.GetConfigFiles(
		ctx.Request.Context(),
		nginx.GetConfigFilesInput{
			BasePath:   normalizePathQuery(ctx, "basePath"),
			ConfigPath: normalizePathQuery(ctx, "configPath"),
			LogPath:    normalizePathQuery(ctx, "logPath"),
			CachePath:  normalizePathQuery(ctx, "cachePath"),
		},
	)
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
