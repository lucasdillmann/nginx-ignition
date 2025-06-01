package access_list

import (
	"dillmann.com.br/nginx-ignition/api/common/authorization"
	"dillmann.com.br/nginx-ignition/core/access_list"
	"dillmann.com.br/nginx-ignition/core/user"
	"github.com/gin-gonic/gin"
)

func Install(router *gin.Engine, commands *access_list.Commands, authorizer *authorization.ABAC) {
	basePath := authorizer.ConfigureGroup(
		router,
		"/api/access-lists",
		func(permissions user.Permissions) user.AccessLevel { return permissions.AccessLists },
	)
	basePath.GET("", listHandler{commands}.handle)
	basePath.POST("", createHandler{commands}.handle)

	byIdPath := basePath.Group("/:id")
	byIdPath.GET("", getHandler{commands}.handle)
	byIdPath.PUT("", updateHandler{commands}.handle)
	byIdPath.DELETE("", deleteHandler{commands}.handle)
}
