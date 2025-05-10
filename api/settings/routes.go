package settings

import (
	"dillmann.com.br/nginx-ignition/api/common/authorization"
	"dillmann.com.br/nginx-ignition/core/settings"
	"dillmann.com.br/nginx-ignition/core/user"
	"github.com/gin-gonic/gin"
)

const (
	apiPath = "/api/settings"
)

func Install(
	router *gin.Engine,
	authorizer *authorization.RBAC,
	commands *settings.Commands,
) {
	basePath := router.Group(apiPath)
	basePath.GET("", getHandler{commands}.handle)
	basePath.PUT("", putHandler{commands}.handle)

	authorizer.RequireRole("PUT", apiPath, user.AdminRole)
}
