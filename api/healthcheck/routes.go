package healthcheck

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"dillmann.com.br/nginx-ignition/api/common/authorization"
	"dillmann.com.br/nginx-ignition/core/common/configuration"
	"dillmann.com.br/nginx-ignition/core/common/healthcheck"
	"dillmann.com.br/nginx-ignition/core/common/log"
)

const (
	apiPath = "/api/health"
)

func Install(
	router *gin.Engine,
	authorizer *authorization.ABAC,
	healthCheck *healthcheck.HealthCheck,
	cfg *configuration.Configuration,
) {
	enabled, err := cfg.GetBoolean("nginx-ignition.health-check.enabled")
	if err != nil {
		log.Warnf("Unable to check if health check endpoints should be enabled or not (error was %v). Keeping them disabled for now.", err)
		return
	}

	if !enabled {
		log.Warnf("Health check endpoints disabled by configuration")
		return
	}

	basePath := router.Group(apiPath)
	basePath.GET("/liveness", livenessHandler{healthCheck}.handle)
	basePath.GET("/readiness", readinessHandler{}.handle)

	authorizer.AllowAnonymous(http.MethodGet, apiPath+"/liveness")
	authorizer.AllowAnonymous(http.MethodGet, apiPath+"/readiness")
}
