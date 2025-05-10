package stream

import (
	"dillmann.com.br/nginx-ignition/core/stream"
	"github.com/gin-gonic/gin"
)

func Install(router *gin.Engine, commands *stream.Commands) {
	basePath := router.Group("/api/streams")
	basePath.GET("", listHandler{commands}.handle)
	basePath.POST("", createHandler{commands}.handle)

	byIdPath := basePath.Group("/:id")
	byIdPath.GET("", getHandler{commands}.handle)
	byIdPath.PUT("", updateHandler{commands}.handle)
	byIdPath.DELETE("", deleteHandler{commands}.handle)
	byIdPath.POST("/toggle-enabled", toggleEnabledHandler{commands}.handle)
}
