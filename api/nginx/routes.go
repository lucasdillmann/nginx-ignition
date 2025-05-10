package nginx

import (
	"dillmann.com.br/nginx-ignition/core/nginx"
	"github.com/gin-gonic/gin"
)

func Install(router *gin.Engine, commands *nginx.Commands) {
	basePath := router.Group("/api/nginx")
	basePath.POST("/start", startHandler{commands}.handle)
	basePath.POST("/stop", stopHandler{commands}.handle)
	basePath.POST("/reload", reloadHandler{commands}.handle)
	basePath.GET("/status", statusHandler{commands}.handle)
	basePath.GET("/logs", logsHandler{commands}.handle)
}
