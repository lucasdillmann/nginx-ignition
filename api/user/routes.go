package user

import (
	"net/http"

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
	totpPath.GET("", totpStatusHandler{commands}.handle)
	totpPath.POST("", totpEnableHandler{commands}.handle)
	totpPath.POST("/activate", totpActivateHandler{commands}.handle)
	totpPath.DELETE("", totpDisableHandler{commands}.handle)

	authorizer.AllowAnonymous(http.MethodGet, "/api/users/onboarding/status")
	authorizer.AllowAnonymous(http.MethodPost, "/api/users/onboarding/finish")
	authorizer.AllowAnonymous(http.MethodPost, "/api/users/login")
	authorizer.AllowAllUsers(http.MethodPost, "/api/users/logout")
	authorizer.AllowAllUsers(http.MethodGet, "/api/users/current")
	authorizer.AllowAllUsers(http.MethodPost, "/api/users/current/update-password")
	authorizer.AllowAllUsers(http.MethodGet, "/api/users/current/totp")
	authorizer.AllowAllUsers(http.MethodPost, "/api/users/current/totp")
	authorizer.AllowAllUsers(http.MethodPost, "/api/users/current/totp/activate")
	authorizer.AllowAllUsers(http.MethodDelete, "/api/users/current/totp")
}
