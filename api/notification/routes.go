package notification

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"dillmann.com.br/nginx-ignition/api/common/authorization"
	"dillmann.com.br/nginx-ignition/core/notification"
)

func Install(
	router *gin.Engine,
	authorizer *authorization.ABAC,
	commands notification.Commands,
) {
	inboxPath := router.Group("/api/notifications")
	inboxPath.GET("/categories", categoriesHandler{commands}.handle)
	inboxPath.GET("/unread-count", unreadCountHandler{commands}.handle)
	inboxPath.POST("/mark-all-as-read", markAllAsReadHandler{commands}.handle)
	inboxPath.GET("", listHandler{commands}.handle)

	byIDPath := inboxPath.Group("/:id")
	byIDPath.GET("", getHandler{commands}.handle)
	byIDPath.POST("/mark-as-read", markAsReadHandler{commands}.handle)

	configPath := router.Group("/api/notifications/configurations")
	configPath.GET("/available-providers", availableProvidersHandler{commands}.handle)
	configPath.GET("", listConfigurationsHandler{commands}.handle)
	configPath.POST("", createConfigurationHandler{commands}.handle)

	configByIDPath := configPath.Group("/:id")
	configByIDPath.GET("", getConfigurationHandler{commands}.handle)
	configByIDPath.PUT("", putConfigurationHandler{commands}.handle)
	configByIDPath.DELETE("", deleteConfigurationHandler{commands}.handle)

	authorizer.AllowAllUsers(http.MethodGet, "/api/notifications/categories")
	authorizer.AllowAllUsers(http.MethodGet, "/api/notifications/unread-count")
	authorizer.AllowAllUsers(http.MethodPost, "/api/notifications/mark-all-as-read")
	authorizer.AllowAllUsers(http.MethodGet, "/api/notifications")
	authorizer.AllowAllUsers(http.MethodGet, "/api/notifications/:id")
	authorizer.AllowAllUsers(http.MethodPost, "/api/notifications/:id/mark-as-read")
	authorizer.AllowAllUsers(
		http.MethodGet,
		"/api/notifications/configurations/available-providers",
	)
	authorizer.AllowAllUsers(http.MethodGet, "/api/notifications/configurations")
	authorizer.AllowAllUsers(http.MethodPost, "/api/notifications/configurations")
	authorizer.AllowAllUsers(http.MethodGet, "/api/notifications/configurations/:id")
	authorizer.AllowAllUsers(http.MethodPut, "/api/notifications/configurations/:id")
	authorizer.AllowAllUsers(http.MethodDelete, "/api/notifications/configurations/:id")
}
