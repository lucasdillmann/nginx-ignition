package integration

import (
	"github.com/gin-gonic/gin"

	"dillmann.com.br/nginx-ignition/api/common/authorization"
	"dillmann.com.br/nginx-ignition/core/integration"
	"dillmann.com.br/nginx-ignition/core/user"
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
	basePath.GET("", listHandler{commands}.handle)
	basePath.POST("", createHandler{commands}.handle)

	basePath.GET("/available-drivers", availableDriversHandler{commands}.handle)

	byIdPath := basePath.Group("/:id")
	byIdPath.GET("", getHandler{commands}.handle)
	byIdPath.PUT("", putHandler{commands}.handle)
	byIdPath.DELETE("", deleteHandler{commands}.handle)

	optionsPath := byIdPath.Group("/options")
	optionsPath.GET("", listOptionsHandler{commands}.handle)
	optionsPath.GET("/:optionId", getOptionHandler{commands}.handle)
}
