package settings

import (
	"github.com/gin-gonic/gin"

	"dillmann.com.br/nginx-ignition/api/common/authorization"
	"dillmann.com.br/nginx-ignition/core/settings"
	"dillmann.com.br/nginx-ignition/core/user"
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
