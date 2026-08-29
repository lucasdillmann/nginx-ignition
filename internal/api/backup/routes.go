package backup

import (
	"github.com/gin-gonic/gin"

	"nginx-ignition/internal/api/common/authorization"
	"nginx-ignition/internal/core/backup"
	"nginx-ignition/internal/core/user"
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
