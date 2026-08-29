package nginx

import (
	"github.com/gin-gonic/gin"

	"nginx-ignition/internal/api/common/authorization"
	"nginx-ignition/internal/core/nginx"
	"nginx-ignition/internal/core/settings"
	"nginx-ignition/internal/core/user"
)

func Install(
	router *gin.Engine,
	nginxCommands nginx.Commands,
	settingsCommands settings.Commands,
	authorizer *authorization.ABAC,
) {
	basePath := authorizer.ConfigureGroup(
		router,
		"/api/nginx",
		func(permissions user.Permissions) user.AccessLevel { return permissions.NginxServer },
	)
	basePath.POST("/start", startHandler{nginxCommands}.handle)
	basePath.POST("/stop", stopHandler{nginxCommands}.handle)
	basePath.POST("/reload", reloadHandler{nginxCommands}.handle)
	basePath.GET("/status", statusHandler{nginxCommands}.handle)
	basePath.GET("/metadata", metadataHandler{nginxCommands, settingsCommands}.handle)

	logsPath := authorizer.ConfigureGroup(
		router,
		"/api/nginx/logs",
		func(permissions user.Permissions) user.AccessLevel { return permissions.Logs },
	)
	logsPath.GET("", logsHandler{nginxCommands}.handle)

	exportPath := authorizer.ConfigureGroup(
		router,
		"/api/nginx/config",
		func(permissions user.Permissions) user.AccessLevel { return permissions.ExportData },
	)
	exportPath.GET("", configFilesHandler{nginxCommands}.handle)

	trafficStatsPath := authorizer.ConfigureGroup(
		router,
		"/api/nginx/traffic-stats",
		func(permissions user.Permissions) user.AccessLevel { return permissions.TrafficStats },
	)
	trafficStatsPath.GET("", trafficStatsHandler{nginxCommands}.handle)
}
