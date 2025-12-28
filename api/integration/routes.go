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

	byIDPath := basePath.Group("/:id")
	byIDPath.GET("", getHandler{commands}.handle)
	byIDPath.PUT("", putHandler{commands}.handle)
	byIDPath.DELETE("", deleteHandler{commands}.handle)

	optionsPath := byIDPath.Group("/options")
	optionsPath.GET("", listOptionsHandler{commands}.handle)
	optionsPath.GET("/:optionID", getOptionHandler{commands}.handle)
}
