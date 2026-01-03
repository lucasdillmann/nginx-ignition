package frontend

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"dillmann.com.br/nginx-ignition/api/common/authorization"
	"dillmann.com.br/nginx-ignition/core/common/configuration"
	"dillmann.com.br/nginx-ignition/core/common/log"
)

func Install(
	router *gin.Engine,
	cfg *configuration.Configuration,
	authorizer *authorization.ABAC,
) {
	settingsHandlerInstance := &configurationHandler{cfg}
	router.Handle(http.MethodGet, "/api/frontend/configuration", settingsHandlerInstance.handle)
	authorizer.AllowAllUsers(http.MethodGet, "/api/frontend/configuration")

	var staticHandler staticFilesHandler
	basePath, err := cfg.Get("nginx-ignition.server.frontend-path")

	if err != nil || basePath == "" {
		log.Warnf(
			"Frontend path is not defined. Every request to it will be rejected with not found status.",
		)
		staticHandler = staticFilesHandler{}
	} else {
		log.Infof("Serving frontend files from %s", basePath)
		staticHandler = staticFilesHandler{&basePath}
	}

	router.NoRoute(staticHandler.handle)
}
