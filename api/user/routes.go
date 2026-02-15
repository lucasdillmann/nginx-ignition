package user

import (
	"github.com/gin-gonic/gin"

	"dillmann.com.br/nginx-ignition/api/common/authorization"
	"dillmann.com.br/nginx-ignition/core/user"
)

func Install(
	router *gin.Engine,
	authorizer *authorization.ABAC,
	commands user.Commands,
) {
	basePath := authorizer.ConfigureGroup(
		router,
		"/api/users",
		func(permissions user.Permissions) user.AccessLevel { return permissions.Users },
	)
	basePath.GET("", listHandler{commands}.handle)
	basePath.POST("", createHandler{commands}.handle)

	byIDPath := basePath.Group("/:id")
	byIDPath.GET("", getHandler{commands}.handle)
	byIDPath.PUT("", updateHandler{commands}.handle)
	byIDPath.DELETE("", deleteHandler{commands}.handle)

	onboardingPath := basePath.Group("/onboarding")
	onboardingPath.GET("/status", onboardingStatusHandler{commands}.handle)
	onboardingPath.POST("/finish", onboardingFinishHandler{commands, authorizer}.handle)

	basePath.POST("/logout", logoutHandler{authorizer}.handle)
	basePath.POST("/login", loginHandler{commands, authorizer}.handle)

	currentPath := basePath.Group("/current")
	currentPath.GET("", currentHandler{}.handle)
	currentPath.POST("/update-password", updatePasswordHandler{commands}.handle)

	totpPath := currentPath.Group("/totp")
	totpPath.GET("/status", totpStatusHandler{commands}.handle)
	totpPath.POST("/enable", totpEnableHandler{commands}.handle)
	totpPath.POST("/activate", totpActivateHandler{commands}.handle)
	totpPath.POST("/disable", totpDisableHandler{commands}.handle)

	authorizer.AllowAnonymous("GET", "/api/users/onboarding/status")
	authorizer.AllowAnonymous("POST", "/api/users/onboarding/finish")
	authorizer.AllowAnonymous("POST", "/api/users/login")
	authorizer.AllowAllUsers("POST", "/api/users/logout")
	authorizer.AllowAllUsers("GET", "/api/users/current")
	authorizer.AllowAllUsers("POST", "/api/users/current/update-password")
	authorizer.AllowAllUsers("GET", "/api/users/current/totp/status")
	authorizer.AllowAllUsers("POST", "/api/users/current/totp/enable")
	authorizer.AllowAllUsers("POST", "/api/users/current/totp/activate")
	authorizer.AllowAllUsers("POST", "/api/users/current/totp/disable")
}
