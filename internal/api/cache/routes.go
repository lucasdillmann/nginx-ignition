package cache

import (
	"github.com/gin-gonic/gin"

	"nginx-ignition/internal/api/common/authorization"
	"nginx-ignition/internal/core/cache"
	"nginx-ignition/internal/core/user"
)

func Install(router *gin.Engine, commands cache.Commands, authorizer *authorization.ABAC) {
	basePath := authorizer.ConfigureGroup(
		router,
		"/api/caches",
		func(permissions user.Permissions) user.AccessLevel { return permissions.Caches },
	)

	basePath.GET("", listHandler{commands}.handle)
	basePath.POST("", createHandler{commands}.handle)

	byIDPath := basePath.Group("/:id")
	byIDPath.GET("", getHandler{commands}.handle)
	byIDPath.PUT("", updateHandler{commands}.handle)
	byIDPath.DELETE("", deleteHandler{commands}.handle)
}
