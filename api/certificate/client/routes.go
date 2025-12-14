package client

import (
	"github.com/gin-gonic/gin"

	"dillmann.com.br/nginx-ignition/api/common/authorization"
	"dillmann.com.br/nginx-ignition/core/certificate/client"
	"dillmann.com.br/nginx-ignition/core/user"
)

func Install(router *gin.Engine, commands *client.Commands, authorizer *authorization.ABAC) {
	basePath := authorizer.ConfigureGroup(
		router,
		"/api/certificates/client",
		func(permissions user.Permissions) user.AccessLevel { return permissions.ClientCertificates },
	)
	basePath.GET("", listHandler{commands}.handle)
	// basePath.POST("", createHandler{commands}.handle)

	byIDPath := basePath.Group("/:id")
	byIDPath.GET("", getHandler{commands}.handle)
	// byIDPath.PUT("", updateHandler{commands}.handle)
	// byIDPath.DELETE("", deleteHandler{commands}.handle)
	// byIDPath.POST("/replace-ca", replaceCAHandler{commands}.handle)
	// byIDPath.POST("/keys", getCAKeysHandler{commands}.handle)
	// byIDPath.POST("/items", createClientHandler{commands}.handle)

	// byItemIDPath := byIDPath.Group("/items/:clientId")
	// byItemIDPath.PUT("", updateClientHandler{commands}.handle)
	// byItemIDPath.DELETE("", deleteClientHandler{commands}.handle)
	// byItemIDPath.GET("/keys", getClientKeysHandler{commands}.handle)

	// TODO: Finish the handlers implementation
}
