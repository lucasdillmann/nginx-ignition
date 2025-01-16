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
	basePath, err := configuration.Get("nginx-ignition.server.frontend-path")
	handlerInstance := handler{&basePath}

	if err != nil {
		log.Warnf("Frontend path is not defined. Every request to it will be rejected with not found status.")
		handlerInstance = handler{}
	}

	router.NoRoute(handlerInstance.handle)
}
