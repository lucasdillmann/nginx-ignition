package certificate

import (
	"dillmann.com.br/nginx-ignition/core/certificate"
	"github.com/gin-gonic/gin"
)

func Install(
	router *gin.Engine,
	deleteCommand certificate.DeleteCommand,
	availableProvidersCommand certificate.AvailableProvidersCommand,
	getCommand certificate.GetCommand,
	listCommand certificate.ListCommand,
	issueCommand certificate.IssueCommand,
	renewCommand certificate.RenewCommand,
) {
	basePath := router.Group("/api/certificates")
	basePath.GET("", listHandler{&listCommand}.handle)
	basePath.POST("/issue", issueHandler{&issueCommand}.handle)
	basePath.GET("/available-providers", availableProvidersHandler{&availableProvidersCommand}.handle)

	byIdPath := basePath.Group("/:id")
	byIdPath.GET("", getHandler{&getCommand}.handle)
	byIdPath.DELETE("", deleteHandler{&deleteCommand}.handle)
	byIdPath.POST("/renew", renewHandler{&renewCommand}.handle)
}
