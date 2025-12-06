package healthcheck

import (
	"github.com/gin-gonic/gin"

	"dillmann.com.br/nginx-ignition/core/common/healthcheck"
)

const (
	apiPath = "/api/health"
)

func Install(
	router *gin.Engine,
	healthCheck *healthcheck.HealthCheck,
) {
	basePath := router.Group(apiPath)
	basePath.GET("/liveness", livenessHandler{healthCheck}.handle)
	basePath.PUT("/readiness", readinessHandler{}.handle)
}
