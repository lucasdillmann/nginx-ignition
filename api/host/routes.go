package host

import (
	"dillmann.com.br/nginx-ignition/core/host"
	"dillmann.com.br/nginx-ignition/core/nginx"
	"github.com/gin-gonic/gin"
)

func Install(router *gin.Engine, hostCommands *host.Commands, nginxCommands *nginx.Commands) {
	basePath := router.Group("/api/hosts")
	basePath.GET("", listHandler{hostCommands}.handle)
	basePath.POST("", createHandler{hostCommands}.handle)

	byIdPath := basePath.Group("/:id")
	byIdPath.GET("", getHandler{hostCommands}.handle)
	byIdPath.PUT("", updateHandler{hostCommands}.handle)
	byIdPath.DELETE("", deleteHandler{hostCommands}.handle)
	byIdPath.POST("/toggle-enabled", toggleEnabledHandler{hostCommands}.handle)
	byIdPath.GET("/logs/:qualifier", logsHandler{nginxCommands}.handle)
}
