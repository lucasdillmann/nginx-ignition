package access_list_api

import (
	"dillmann.com.br/nginx-ignition/core/access_list"
	"github.com/gin-gonic/gin"
)

const (
	basePath = "/api/access-lists"
	byIdPath = basePath + "/:id"
)

func Install(
	engine *gin.Engine,
	getCommand access_list.GetCommand,
	saveCommand access_list.SaveCommand,
	deleteCommand access_list.DeleteCommand,
	listCommand access_list.ListCommand,
) {
	engine.GET(basePath, listHandler{&listCommand}.handle)
	engine.POST(basePath, createHandler{&saveCommand}.handle)
	engine.GET(byIdPath, getHandler{&getCommand}.handle)
	engine.PUT(byIdPath, updateHandler{&saveCommand}.handle)
	engine.DELETE(byIdPath, deleteHandler{&deleteCommand}.handle)
}
