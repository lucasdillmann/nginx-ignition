package backup

import (
	"github.com/gin-gonic/gin"

	"dillmann.com.br/nginx-ignition/api/common/authorization"
	"dillmann.com.br/nginx-ignition/core/backup"
	"dillmann.com.br/nginx-ignition/core/user"
)

const (
	apiPath = "/api/backup"
)

func Install(
	router *gin.Engine,
	authorizer *authorization.ABAC,
	commands backup.Commands,
) {
	basePath := authorizer.ConfigureGroup(
		router,
		apiPath,
		func(permissions user.Permissions) user.AccessLevel { return permissions.ExportData },
	)
	basePath.GET("", getHandler{commands}.handle)
}
