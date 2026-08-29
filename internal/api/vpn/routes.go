package vpn

import (
	"github.com/gin-gonic/gin"

	"nginx-ignition/internal/api/common/authorization"
	"nginx-ignition/internal/core/user"
	"nginx-ignition/internal/core/vpn"
)

func Install(
	router *gin.Engine,
	authorizer *authorization.ABAC,
	commands vpn.Commands,
) {
	basePath := authorizer.ConfigureGroup(
		router,
		"/api/vpns",
		func(permissions user.Permissions) user.AccessLevel { return permissions.VPNs },
	)
	basePath.GET("", listHandler{commands}.handle)
	basePath.POST("", createHandler{commands}.handle)

	basePath.GET("/available-drivers", availableDriversHandler{commands}.handle)

	byIDPath := basePath.Group("/:id")
	byIDPath.GET("", getHandler{commands}.handle)
	byIDPath.PUT("", putHandler{commands}.handle)
	byIDPath.DELETE("", deleteHandler{commands}.handle)
}
