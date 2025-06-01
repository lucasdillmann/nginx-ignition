package nginx

import (
	"dillmann.com.br/nginx-ignition/api/common/authorization"
	"dillmann.com.br/nginx-ignition/core/nginx"
	"dillmann.com.br/nginx-ignition/core/user"
	"github.com/gin-gonic/gin"
)

func Install(router *gin.Engine, commands *nginx.Commands, authorizer *authorization.ABAC) {
	basePath := authorizer.ConfigureGroup(
		router,
		"/api/nginx",
		func(permissions user.Permissions) user.AccessLevel { return permissions.NginxServer },
	)
	basePath.POST("/start", startHandler{commands}.handle)
	basePath.POST("/stop", stopHandler{commands}.handle)
	basePath.POST("/reload", reloadHandler{commands}.handle)
	basePath.GET("/status", statusHandler{commands}.handle)

	logsPath := authorizer.ConfigureGroup(
		router,
		"/api/nginx/logs",
		func(permissions user.Permissions) user.AccessLevel { return permissions.Logs },
	)
	logsPath.GET("", logsHandler{commands}.handle)
}
