package access_list

import (
	"dillmann.com.br/nginx-ignition/core/access_list"
	"github.com/gin-gonic/gin"
)

const (
	basePath = "/api/access-lists"
	byIdPath = basePath + "/:id"
)

func Install(
	router *gin.Engine,
	getCommand access_list.GetCommand,
	saveCommand access_list.SaveCommand,
	deleteCommand access_list.DeleteCommand,
	listCommand access_list.ListCommand,
) {
	router.GET(basePath, listHandler{&listCommand}.handle)
	router.POST(basePath, createHandler{&saveCommand}.handle)

	router.GET(byIdPath, getHandler{&getCommand}.handle)
	router.PUT(byIdPath, updateHandler{&saveCommand}.handle)
	router.DELETE(byIdPath, deleteHandler{&deleteCommand}.handle)
}
