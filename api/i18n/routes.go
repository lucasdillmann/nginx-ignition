package i18n

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"dillmann.com.br/nginx-ignition/api/common/authorization"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
)

func Install(router *gin.Engine, commands i18n.Commands, authorizer *authorization.ABAC) {
	router.Use(middleware(commands))

	basePath := router.Group("/api/i18n")
	basePath.GET("", getDictionaryHandler{commands}.handle)

	authorizer.AllowAnonymous(http.MethodGet, "/api/i18n")
}
