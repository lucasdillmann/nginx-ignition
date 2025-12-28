package stream

import (
	"github.com/gin-gonic/gin"

	"dillmann.com.br/nginx-ignition/api/common/authorization"
	"dillmann.com.br/nginx-ignition/core/stream"
	"dillmann.com.br/nginx-ignition/core/user"
)

func Install(router *gin.Engine, commands *stream.Commands, authorizer *authorization.ABAC) {
	basePath := authorizer.ConfigureGroup(
		router,
		"/api/streams",
		func(permissions user.Permissions) user.AccessLevel { return permissions.Streams },
	)
	basePath.GET("", listHandler{commands}.handle)
	basePath.POST("", createHandler{commands}.handle)

	byIDPath := basePath.Group("/:id")
	byIDPath.GET("", getHandler{commands}.handle)
	byIDPath.PUT("", updateHandler{commands}.handle)
	byIDPath.DELETE("", deleteHandler{commands}.handle)
	byIDPath.POST("/toggle-enabled", toggleEnabledHandler{commands}.handle)
}
