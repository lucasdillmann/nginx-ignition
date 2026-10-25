package certificate

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/lucasdillmann/nginx-ignition/internal/api/core/authorization"
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/log"

	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/certificate"
	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/user"
)

const extendedWriteDeadline = 10 * time.Minute

func Install(router *gin.Engine, commands certificate.Commands, authorizer *authorization.ABAC) {
	basePath := authorizer.ConfigureGroup(
		router,
		"/api/certificates",
		func(permissions user.Permissions) user.AccessLevel { return permissions.Certificates },
	)
	basePath.GET("", listHandler{commands}.handle)
	basePath.POST("/issue", extendWriteDeadline, issueHandler{commands}.handle)
	basePath.GET("/available-providers", availableProvidersHandler{commands}.handle)

	byIDPath := basePath.Group("/:id")
	byIDPath.GET("", getHandler{commands}.handle)
	byIDPath.DELETE("", deleteHandler{commands}.handle)
	byIDPath.POST("/renew", extendWriteDeadline, renewHandler{commands}.handle)
}

func extendWriteDeadline(ctx *gin.Context) {
	controller := http.NewResponseController(ctx.Writer)
	if err := controller.SetWriteDeadline(time.Now().Add(extendedWriteDeadline)); err != nil {
		log.Warnf("Unable to extend the certificate endpoint write deadline: %s", err)
	}

	ctx.Next()
}
