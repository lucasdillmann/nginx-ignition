package accesslist

import (
	"github.com/gin-gonic/gin"

	"github.com/lucasdillmann/nginx-ignition/internal/api/common/authorization"
	"github.com/lucasdillmann/nginx-ignition/internal/core/accesslist"
	"github.com/lucasdillmann/nginx-ignition/internal/core/user"
)

func Install(router *gin.Engine, commands accesslist.Commands, authorizer *authorization.ABAC) {
	basePath := authorizer.ConfigureGroup(
		router,
		"/api/access-lists",
		func(permissions user.Permissions) user.AccessLevel { return permissions.AccessLists },
	)
	basePath.GET("", listHandler{commands}.handle)
	basePath.POST("", createHandler{commands}.handle)

	byIDPath := basePath.Group("/:id")
	byIDPath.GET("", getHandler{commands}.handle)
	byIDPath.PUT("", updateHandler{commands}.handle)
	byIDPath.DELETE("", deleteHandler{commands}.handle)
}
