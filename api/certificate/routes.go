package certificate

import (
	"dillmann.com.br/nginx-ignition/core/certificate"
	"github.com/gin-gonic/gin"
)

func Install(router *gin.Engine, commands *certificate.Commands) {
	basePath := router.Group("/api/certificates")
	basePath.GET("", listHandler{commands}.handle)
	basePath.POST("/issue", issueHandler{commands}.handle)
	basePath.GET("/available-providers", availableProvidersHandler{commands}.handle)

	byIdPath := basePath.Group("/:id")
	byIdPath.GET("", getHandler{commands}.handle)
	byIdPath.DELETE("", deleteHandler{commands}.handle)
	byIdPath.POST("/renew", renewHandler{commands}.handle)
}
