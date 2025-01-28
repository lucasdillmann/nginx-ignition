package integration

import (
	"dillmann.com.br/nginx-ignition/api/common/authorization"
	"dillmann.com.br/nginx-ignition/core/integration"
	"dillmann.com.br/nginx-ignition/core/user"
	"github.com/gin-gonic/gin"
)

func Install(
	router *gin.Engine,
	authorizer *authorization.RBAC,
	listCommand integration.ListCommand,
	getCommand integration.GetByIdCommand,
	setConfigurationCommand integration.ConfigureByIdCommand,
	listOptionsCommand integration.ListOptionsCommand,
	getOptionCommand integration.GetOptionByIdCommand,
) {
	basePath := router.Group("/api/integrations")
	basePath.GET("", listIntegrationsHandler{&listCommand}.handle)

	byIdPath := basePath.Group("/:id")

	optionsPath := byIdPath.Group("/options")
	optionsPath.GET("", listOptionsHandler{&listOptionsCommand}.handle)
	optionsPath.GET("/:optionId", getOptionHandler{&getOptionCommand}.handle)

	configurationPath := byIdPath.Group("/configuration")
	configurationPath.GET("", getConfigurationHandler{&getCommand}.handle)
	configurationPath.PUT("", putConfigurationHandler{&setConfigurationCommand}.handle)

	authorizer.RequireRole("GET", "/api/integrations/:id/configuration", user.AdminRole)
	authorizer.RequireRole("PUT", "/api/integrations/:id/configuration", user.AdminRole)
}
