package i18n

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"nginx-ignition/internal/api/common/authorization"
	"nginx-ignition/internal/core/common/i18n"
)

func Install(router *gin.Engine, commands i18n.Commands, authorizer *authorization.ABAC) {
	basePath := router.Group("/api/i18n")
	basePath.GET("", getAvailableLanguagesHandler{commands}.handle)
	basePath.GET("/:language", getDictionaryHandler{commands}.handle)

	authorizer.AllowAnonymous(http.MethodGet, "/api/i18n")
	authorizer.AllowAnonymous(http.MethodGet, "/api/i18n/:language")
}
