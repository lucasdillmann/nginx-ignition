package access_list

import (
	"dillmann.com.br/nginx-ignition/core/access_list"
	"github.com/gin-gonic/gin"
)

func Install(
	router *gin.Engine,
	getCommand access_list.GetCommand,
	saveCommand access_list.SaveCommand,
	deleteCommand access_list.DeleteCommand,
	listCommand access_list.ListCommand,
) {
	basePath := router.Group("/api/access-lists")
	basePath.GET("", listHandler{&listCommand}.handle)
	basePath.POST("", createHandler{&saveCommand}.handle)

	byIdPath := basePath.Group("/:id")
	byIdPath.GET("", getHandler{&getCommand}.handle)
	byIdPath.PUT("", updateHandler{&saveCommand}.handle)
	byIdPath.DELETE("", deleteHandler{&deleteCommand}.handle)
}
