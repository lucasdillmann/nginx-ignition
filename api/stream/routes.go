package stream

import (
	"dillmann.com.br/nginx-ignition/api/common/authorization"
	"dillmann.com.br/nginx-ignition/core/stream"
	"dillmann.com.br/nginx-ignition/core/user"
	"github.com/gin-gonic/gin"
)

func Install(router *gin.Engine, commands *stream.Commands, authorizer *authorization.ABAC) {
	basePath := authorizer.ConfigureGroup(
		router,
		"/api/streams",
		func(permissions user.Permissions) user.AccessLevel { return permissions.Streams },
	)
	basePath.GET("", listHandler{commands}.handle)
	basePath.POST("", createHandler{commands}.handle)

	byIdPath := basePath.Group("/:id")
	byIdPath.GET("", getHandler{commands}.handle)
	byIdPath.PUT("", updateHandler{commands}.handle)
	byIdPath.DELETE("", deleteHandler{commands}.handle)
	byIdPath.POST("/toggle-enabled", toggleEnabledHandler{commands}.handle)
}
