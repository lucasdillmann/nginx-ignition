package settings_api

import (
	"dillmann.com.br/nginx-ignition/api/common/authentication"
	"dillmann.com.br/nginx-ignition/core/settings"
	"dillmann.com.br/nginx-ignition/core/user"
	"github.com/gin-gonic/gin"
)

const (
	basePath = "/api/settings"
)

func Install(
	engine *gin.Engine,
	authorizer *authentication.RBAC,
	getCommand settings.GetCommand,
	saveCommand settings.SaveCommand,
) {
	engine.GET(basePath, getHandler{&getCommand}.handle)
	engine.PUT(basePath, putHandler{&saveCommand}.handle)

	authorizer.RequireRole("PUT", basePath, user.AdminRole)
}
