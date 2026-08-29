package settings

import (
	"github.com/gin-gonic/gin"

	"nginx-ignition/internal/api/common/authorization"
	"nginx-ignition/internal/core/settings"
	"nginx-ignition/internal/core/user"
)

const (
	apiPath = "/api/settings"
)

func Install(
	router *gin.Engine,
	authorizer *authorization.ABAC,
	commands settings.Commands,
) {
	basePath := authorizer.ConfigureGroup(
		router,
		apiPath,
		func(permissions user.Permissions) user.AccessLevel { return permissions.Settings },
	)
	basePath.GET("", getHandler{commands}.handle)
	basePath.PUT("", putHandler{commands}.handle)
}
