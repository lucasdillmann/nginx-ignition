package server

import (
	"github.com/gin-gonic/gin"

	"dillmann.com.br/nginx-ignition/api/common/authorization"
	"dillmann.com.br/nginx-ignition/core/certificate/server"
	"dillmann.com.br/nginx-ignition/core/user"
)

func Install(router *gin.Engine, commands *server.Commands, authorizer *authorization.ABAC) {
	basePath := authorizer.ConfigureGroup(
		router,
		"/api/certificates/server",
		func(permissions user.Permissions) user.AccessLevel { return permissions.ServerCertificates },
	)
	basePath.GET("", listHandler{commands}.handle)
	basePath.POST("/issue", issueHandler{commands}.handle)
	basePath.GET("/available-providers", availableProvidersHandler{commands}.handle)

	byIdPath := basePath.Group("/:id")
	byIdPath.GET("", getHandler{commands}.handle)
	byIdPath.DELETE("", deleteHandler{commands}.handle)
	byIdPath.POST("/renew", renewHandler{commands}.handle)
}
