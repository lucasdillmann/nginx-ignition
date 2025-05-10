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
	commands *integration.Commands,
) {
	basePath := router.Group("/api/integrations")
	basePath.GET("", listIntegrationsHandler{commands}.handle)

	byIdPath := basePath.Group("/:id")

	optionsPath := byIdPath.Group("/options")
	optionsPath.GET("", listOptionsHandler{commands}.handle)
	optionsPath.GET("/:optionId", getOptionHandler{commands}.handle)

	configurationPath := byIdPath.Group("/configuration")
	configurationPath.GET("", getConfigurationHandler{commands}.handle)
	configurationPath.PUT("", putConfigurationHandler{commands}.handle)

	authorizer.RequireRole("GET", "/api/integrations/:id/configuration", user.AdminRole)
	authorizer.RequireRole("PUT", "/api/integrations/:id/configuration", user.AdminRole)
}
