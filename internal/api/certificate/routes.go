package certificate

import (
	"github.com/gin-gonic/gin"

	"github.com/lucasdillmann/nginx-ignition/internal/api/common/authorization"
	"github.com/lucasdillmann/nginx-ignition/internal/core/certificate"
	"github.com/lucasdillmann/nginx-ignition/internal/core/user"
)

func Install(router *gin.Engine, commands certificate.Commands, authorizer *authorization.ABAC) {
	basePath := authorizer.ConfigureGroup(
		router,
		"/api/certificates",
		func(permissions user.Permissions) user.AccessLevel { return permissions.Certificates },
	)
	basePath.GET("", listHandler{commands}.handle)
	basePath.POST("/issue", issueHandler{commands}.handle)
	basePath.GET("/available-providers", availableProvidersHandler{commands}.handle)

	byIDPath := basePath.Group("/:id")
	byIDPath.GET("", getHandler{commands}.handle)
	byIDPath.DELETE("", deleteHandler{commands}.handle)
	byIDPath.POST("/renew", renewHandler{commands}.handle)
}
