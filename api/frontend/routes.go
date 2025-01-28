package frontend

import (
	"dillmann.com.br/nginx-ignition/core/common/configuration"
	"dillmann.com.br/nginx-ignition/core/common/log"
	"github.com/gin-gonic/gin"
)

func Install(
	router *gin.Engine,
	configuration *configuration.Configuration,
) {
	var handlerInstance handler
	basePath, err := configuration.Get("nginx-ignition.server.frontend-path")

	if err != nil || basePath == "" {
		log.Warnf("Frontend path is not defined. Every request to it will be rejected with not found status.")
		handlerInstance = handler{}
	} else {
		log.Infof("Serving frontend files from %s", basePath)
		handlerInstance = handler{&basePath}
	}

	router.NoRoute(handlerInstance.handle)
}
