package host

import (
	"dillmann.com.br/nginx-ignition/core/host"
	"dillmann.com.br/nginx-ignition/core/nginx"
	"github.com/gin-gonic/gin"
)

func Install(
	router *gin.Engine,
	getCommand host.GetCommand,
	saveCommand host.SaveCommand,
	deleteCommand host.DeleteCommand,
	listCommand host.ListCommand,
	logsCommand nginx.GetHostLogsCommand,
) {
	basePath := router.Group("/api/hosts")
	basePath.GET("", listHandler{&listCommand}.handle)
	basePath.POST("", createHandler{&saveCommand}.handle)

	byIdPath := basePath.Group("/:id")
	byIdPath.GET("", getHandler{&getCommand}.handle)
	byIdPath.PUT("", updateHandler{&saveCommand}.handle)
	byIdPath.DELETE("", deleteHandler{&deleteCommand}.handle)
	byIdPath.POST("/toggle-enabled", toggleEnabledHandler{&getCommand, &saveCommand}.handle)
	byIdPath.GET("/logs/:qualifier", logsHandler{&logsCommand}.handle)
}
