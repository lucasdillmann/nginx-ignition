package access_list

import (
	"dillmann.com.br/nginx-ignition/core/access_list"
	"github.com/gin-gonic/gin"
)

func Install(router *gin.Engine, commands *access_list.Commands) {
	basePath := router.Group("/api/access-lists")
	basePath.GET("", listHandler{commands}.handle)
	basePath.POST("", createHandler{commands}.handle)

	byIdPath := basePath.Group("/:id")
	byIdPath.GET("", getHandler{commands}.handle)
	byIdPath.PUT("", updateHandler{commands}.handle)
	byIdPath.DELETE("", deleteHandler{commands}.handle)
}
