package host

import (
	"dillmann.com.br/nginx-ignition/core/host"
	"github.com/gin-gonic/gin"
)

const (
	basePath = "/api/hosts"
	byIdPath = basePath + "/:id"
)

func Install(
	router *gin.Engine,
	getCommand host.GetCommand,
	saveCommand host.SaveCommand,
	deleteCommand host.DeleteCommand,
	listCommand host.ListCommand,
) {
	router.GET(basePath, listHandler{&listCommand}.handle)
	router.POST(basePath, createHandler{&saveCommand}.handle)

	router.GET(byIdPath, getHandler{&getCommand}.handle)
	router.PUT(byIdPath, updateHandler{&saveCommand}.handle)
	router.DELETE(byIdPath, deleteHandler{&deleteCommand}.handle)
}
