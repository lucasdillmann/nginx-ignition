package stream

import (
	"github.com/gin-gonic/gin"

	"github.com/lucasdillmann/nginx-ignition/internal/api/common/authorization"
	"github.com/lucasdillmann/nginx-ignition/internal/core/stream"
	"github.com/lucasdillmann/nginx-ignition/internal/core/user"
)

func Install(router *gin.Engine, commands stream.Commands, authorizer *authorization.ABAC) {
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
