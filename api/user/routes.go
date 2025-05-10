package user

import (
	"dillmann.com.br/nginx-ignition/api/common/authorization"
	"dillmann.com.br/nginx-ignition/core/user"
	"github.com/gin-gonic/gin"
)

const (
	basePath             = "/api/users"
	byIdPath             = basePath + "/:id"
	onboardingPath       = basePath + "/onboarding"
	onboardingStatusPath = onboardingPath + "/status"
	onboardingFinishPath = onboardingPath + "/finish"
	logoutPath           = basePath + "/logout"
	loginPath            = basePath + "/login"
	currentPath          = basePath + "/current"
	updatePasswordPath   = currentPath + "/update-password"
)

func Install(
	router *gin.Engine,
	authorizer *authorization.RBAC,
	commands *user.Commands,
) {
	router.GET(basePath, listHandler{commands}.handle)
	router.POST(basePath, createHandler{commands}.handle)

	router.GET(byIdPath, getHandler{commands}.handle)
	router.PUT(byIdPath, updateHandler{commands}.handle)
	router.DELETE(byIdPath, deleteHandler{commands}.handle)

	router.GET(onboardingStatusPath, onboardingStatusHandler{commands}.handle)
	router.POST(onboardingFinishPath, onboardingFinishHandler{commands, authorizer}.handle)

	router.POST(logoutPath, logoutHandler{authorizer}.handle)
	router.POST(loginPath, loginHandler{commands, authorizer}.handle)

	router.GET(currentPath, currentHandler{}.handle)
	router.POST(updatePasswordPath, updatePasswordHandler{commands}.handle)

	authorizer.RequireRole("GET", basePath, user.AdminRole)
	authorizer.RequireRole("POST", basePath, user.AdminRole)

	authorizer.RequireRole("GET", byIdPath, user.AdminRole)
	authorizer.RequireRole("PUT", byIdPath, user.AdminRole)
	authorizer.RequireRole("DELETE", byIdPath, user.AdminRole)

	authorizer.AllowAnonymous("GET", onboardingStatusPath)
	authorizer.AllowAnonymous("POST", onboardingFinishPath)

	authorizer.AllowAnonymous("POST", loginPath)
}
