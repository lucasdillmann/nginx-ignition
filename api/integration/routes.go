package integration

import (
	"dillmann.com.br/nginx-ignition/api/common/authorization"
	"dillmann.com.br/nginx-ignition/core/integration"
	"dillmann.com.br/nginx-ignition/core/user"
	"github.com/gin-gonic/gin"
)

func Install(
	router *gin.Engine,
	authorizer *authorization.ABAC,
	commands *integration.Commands,
) {
	basePath := authorizer.ConfigureGroup(
		router,
		"/api/integrations",
		func(permissions user.Permissions) user.AccessLevel { return permissions.Integrations },
	)
	basePath.GET("", listIntegrationsHandler{commands}.handle)

	byIdPath := basePath.Group("/:id")
	optionsPath := byIdPath.Group("/options")
	optionsPath.GET("", listOptionsHandler{commands}.handle)
	optionsPath.GET("/:optionId", getOptionHandler{commands}.handle)

	configurationPath := byIdPath.Group("/configuration")
	configurationPath.GET("", getConfigurationHandler{commands}.handle)
	configurationPath.PUT("", putConfigurationHandler{commands}.handle)
}
