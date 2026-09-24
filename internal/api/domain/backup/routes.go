package backup

import (
	"github.com/gin-gonic/gin"

	"github.com/lucasdillmann/nginx-ignition/internal/api/core/authorization"

	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/backup"
	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/user"
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
