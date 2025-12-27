package cache

import (
	"github.com/gin-gonic/gin"

	"dillmann.com.br/nginx-ignition/api/common/authorization"
	"dillmann.com.br/nginx-ignition/core/cache"
	"dillmann.com.br/nginx-ignition/core/user"
)

func Install(router *gin.Engine, commands *cache.Commands, authorizer *authorization.ABAC) {
	basePath := authorizer.ConfigureGroup(
		router,
		"/api/caches",
		func(permissions user.Permissions) user.AccessLevel { return permissions.Caches },
	)

	basePath.GET("", listHandler{commands}.handle)
	basePath.POST("", createHandler{commands}.handle)

	byIdPath := basePath.Group("/:id")
	byIdPath.GET("", getHandler{commands}.handle)
	byIdPath.PUT("", updateHandler{commands}.handle)
	byIdPath.DELETE("", deleteHandler{commands}.handle)
}
