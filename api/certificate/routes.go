package certificate

import (
	"github.com/gin-gonic/gin"

	"dillmann.com.br/nginx-ignition/api/common/authorization"
	"dillmann.com.br/nginx-ignition/core/certificate"
	"dillmann.com.br/nginx-ignition/core/user"
)

func Install(router *gin.Engine, commands *certificate.Commands, authorizer *authorization.ABAC) {
	basePath := authorizer.ConfigureGroup(
		router,
		"/api/certificates",
		func(permissions user.Permissions) user.AccessLevel { return permissions.Certificates },
	)
	basePath.GET("", listHandler{commands}.handle)
	basePath.POST("/issue", issueHandler{commands}.handle)
	basePath.GET("/available-providers", availableProvidersHandler{commands}.handle)

	byIdPath := basePath.Group("/:id")
	byIdPath.GET("", getHandler{commands}.handle)
	byIdPath.DELETE("", deleteHandler{commands}.handle)
	byIdPath.POST("/renew", renewHandler{commands}.handle)
}
