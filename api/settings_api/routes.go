package settings_api

import (
	"dillmann.com.br/nginx-ignition/core/settings"
	"github.com/gin-gonic/gin"
)

func Install(
	engine *gin.Engine,
	getCommand settings.GetCommand,
) {
	engine.GET("/api/settings", getHandler{getCommand}.handle)
}
