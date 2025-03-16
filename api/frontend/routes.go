package frontend

import (
	"dillmann.com.br/nginx-ignition/core/common/configuration"
	"dillmann.com.br/nginx-ignition/core/common/log"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Install(
	router *gin.Engine,
	configuration *configuration.Configuration,
) {
	settingsHandlerInstance := &configurationHandler{configuration}
	router.Handle(http.MethodGet, "/api/frontend/configuration", settingsHandlerInstance.handle)

	var staticHandler staticFilesHandler
	basePath, err := configuration.Get("nginx-ignition.server.frontend-path")

	if err != nil || basePath == "" {
		log.Warnf("Frontend path is not defined. Every request to it will be rejected with not found status.")
		staticHandler = staticFilesHandler{}
	} else {
		log.Infof("Serving frontend files from %s", basePath)
		staticHandler = staticFilesHandler{&basePath}
	}

	router.NoRoute(staticHandler.handle)
}
