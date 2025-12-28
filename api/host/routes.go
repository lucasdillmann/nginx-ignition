package host

import (
	"github.com/gin-gonic/gin"

	"dillmann.com.br/nginx-ignition/api/common/authorization"
	"dillmann.com.br/nginx-ignition/core/host"
	"dillmann.com.br/nginx-ignition/core/nginx"
	"dillmann.com.br/nginx-ignition/core/settings"
	"dillmann.com.br/nginx-ignition/core/user"
)

func Install(
	router *gin.Engine,
	hostCommands *host.Commands,
	nginxCommands *nginx.Commands,
	settingsCommands *settings.Commands,
	authorizer *authorization.ABAC,
) {
	basePath := authorizer.ConfigureGroup(
		router,
		"/api/hosts",
		func(permissions user.Permissions) user.AccessLevel { return permissions.Hosts },
	)
	basePath.GET("", listHandler{settingsCommands, hostCommands}.handle)
	basePath.POST("", createHandler{hostCommands}.handle)

	byIDPath := basePath.Group("/:id")
	byIDPath.GET("", getHandler{settingsCommands, hostCommands}.handle)
	byIDPath.PUT("", updateHandler{hostCommands}.handle)
	byIDPath.DELETE("", deleteHandler{hostCommands}.handle)
	byIDPath.POST("/toggle-enabled", toggleEnabledHandler{hostCommands}.handle)

	logsPath := authorizer.ConfigureGroup(
		router,
		"/api/hosts/:id/logs",
		func(permissions user.Permissions) user.AccessLevel { return permissions.Logs },
	)
	logsPath.GET("/:qualifier", logsHandler{nginxCommands}.handle)
}
