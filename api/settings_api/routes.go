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
	router *gin.Engine,
	authorizer *authentication.RBAC,
	getCommand settings.GetCommand,
	saveCommand settings.SaveCommand,
) {
	router.GET(basePath, getHandler{&getCommand}.handle)
	router.PUT(basePath, putHandler{&saveCommand}.handle)

	authorizer.RequireRole("PUT", basePath, user.AdminRole)
}
