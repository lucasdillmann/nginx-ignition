package user_api

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
	getCommand user.GetCommand,
	saveCommand user.SaveCommand,
	deleteCommand user.DeleteCommand,
	listCommand user.ListCommand,
	updatePasswordCommand user.UpdatePasswordCommand,
	authenticateCommand user.AuthenticateCommand,
	onboardingStatusCommand user.OnboardingCompletedCommand,
) {
	router.GET(basePath, listHandler{&listCommand}.handle)
	router.POST(basePath, createHandler{&saveCommand}.handle)

	router.GET(byIdPath, getHandler{&getCommand}.handle)
	router.PUT(byIdPath, updateHandler{&saveCommand}.handle)
	router.DELETE(byIdPath, deleteHandler{&deleteCommand}.handle)

	router.GET(onboardingStatusPath, onboardingStatusHandler{&onboardingStatusCommand}.handle)
	router.POST(onboardingFinishPath, onboardingFinishHandler{
		&onboardingStatusCommand,
		&saveCommand,
		&authenticateCommand,
		authorizer,
	}.handle)

	router.POST(logoutPath, logoutHandler{authorizer}.handle)
	router.POST(loginPath, loginHandler{&authenticateCommand, authorizer}.handle)

	router.GET(currentPath, currentHandler{}.handle)
	router.POST(updatePasswordPath, updatePasswordHandler{&updatePasswordCommand}.handle)

	authorizer.RequireRole("GET", basePath, user.AdminRole)
	authorizer.RequireRole("POST", basePath, user.AdminRole)

	authorizer.RequireRole("GET", byIdPath, user.AdminRole)
	authorizer.RequireRole("PUT", byIdPath, user.AdminRole)
	authorizer.RequireRole("DELETE", byIdPath, user.AdminRole)

	authorizer.AllowAnonymous("GET", onboardingStatusPath)
	authorizer.AllowAnonymous("POST", onboardingFinishPath)

	authorizer.AllowAnonymous("POST", loginPath)
}
