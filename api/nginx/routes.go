package nginx

import (
	"dillmann.com.br/nginx-ignition/core/nginx"
	"github.com/gin-gonic/gin"
)

func Install(
	router *gin.Engine,
	reloadCommand nginx.ReloadCommand,
	startCommand nginx.StartCommand,
	stopCommand nginx.StopCommand,
	statusCommand nginx.GetStatusCommand,
	logsCommand nginx.GetMainLogsCommand,
) {
	basePath := router.Group("/api/nginx")
	basePath.POST("/start", startHandler{&startCommand}.handle)
	basePath.POST("/stop", stopHandler{&stopCommand}.handle)
	basePath.POST("/reload", reloadHandler{&reloadCommand}.handle)
	basePath.GET("/status", statusHandler{&statusCommand}.handle)
	basePath.GET("/logs", logsHandler{&logsCommand}.handle)
}
