package host

import (
	"github.com/gin-gonic/gin"

	"github.com/lucasdillmann/nginx-ignition/internal/api/common/authorization"
	"github.com/lucasdillmann/nginx-ignition/internal/core/host"
	"github.com/lucasdillmann/nginx-ignition/internal/core/nginx"
	"github.com/lucasdillmann/nginx-ignition/internal/core/settings"
	"github.com/lucasdillmann/nginx-ignition/internal/core/user"
)

func Install(
	router *gin.Engine,
	hostCommands host.Commands,
	nginxCommands nginx.Commands,
	settingsCommands settings.Commands,
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
